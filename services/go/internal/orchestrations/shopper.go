package orchestrations

import (
	"errors"
	"fmt"
	commands "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/commands/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/orchestrations/v1"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	inventory2 "github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"github.com/temporalio/temporal-shop/services/go/internal/shopping"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
	"time"
)

const SessionMinSeconds = 3600 * 24 * 60
const ShopperRefreshCountThreshold = 500

type ShopperState struct {
	InventoryID string
	ShopperID   string
	Email       string
}

// ExpiringSession is an example of a timer that can have its duration updated
type ExpiringSession struct {
	params            *orchestrations2.StartShopperRequest
	durationSeconds   time.Duration
	refreshesReceived int
}

// SleepUntil sleeps until the provided wake-up time.
// The wake-up time can be updated at any time by sending a new time over updateWakeUpTimeCh.
// Supports ctx cancellation.
// Returns temporal.CanceledError if ctx was canceled.
// Returns temporal.ContinueAsNewError if threshold for refreshes is met
func (u *ExpiringSession) SleepUntil(
	ctx workflow.Context,
	durationSeconds time.Duration,
	refreshShopperChan workflow.ReceiveChannel,
) (err error) {
	logger := workflow.GetLogger(ctx)
	u.durationSeconds = durationSeconds
	timerFired := false
	for !timerFired && ctx.Err() == nil {
		timerCtx, timerCancel := workflow.WithCancel(ctx)
		timer := workflow.NewTimer(timerCtx, u.durationSeconds)
		logger.Info("SleepUntil", "duration_seconds", u.durationSeconds)
		workflow.NewSelector(timerCtx).
			AddFuture(timer, func(f workflow.Future) {
				err := f.Get(timerCtx, nil)
				// if a timer returned an error then it was canceled
				if err == nil {
					logger.Info("Timer fired")
					timerFired = true
				} else if ctx.Err() != nil { // Only log on root ctx cancellation, not on timerCancel function call.
					logger.Info("SleepUntil canceled")
				}
			}).
			AddReceive(refreshShopperChan, func(c workflow.ReceiveChannel, more bool) {
				timerCancel() // cancel outstanding timer
				req := &commands.RefreshShopperRequest{}
				c.Receive(timerCtx, &req) // update wake-up time
				if req.DurationSeconds == 0 {
					req.DurationSeconds = SessionMinSeconds
				}
				u.refreshesReceived = u.refreshesReceived + 1
				u.durationSeconds = time.Second * time.Duration(req.DurationSeconds)
				logger.Info("Wake up time update requested")
			}).
			Select(timerCtx)
		// we've received many events so wipe the slate clean
		if u.refreshesReceived >= ShopperRefreshCountThreshold {
			return workflow.NewContinueAsNewError(ctx, TypeOrchestrations.Shopper, u.params)
		}
	}
	return ctx.Err()
}

func (w *Orchestrations) Shopper(ctx workflow.Context, params *orchestrations2.StartShopperRequest) error {
	if params == nil {
		return fmt.Errorf("params are required")
	}
	/* TODO explicit validation function */
	if params.DurationSeconds == 0 {
		params.DurationSeconds = SessionMinSeconds
	}
	if params.InventoryId == "" {
		params.InventoryId = inventory2.InventorySessionID(params.ShopperId)
	}
	if params.CartId == "" {
		params.CartId = shopping.CartID(params.ShopperId)
	}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 3,
	})
	logger := log.With(workflow.GetLogger(ctx), "email", params.Email)
	logger.Info("started shopper")

	state := &queries.GetShopperResponse{
		ShopperId:   params.ShopperId,
		Email:       params.Email,
		InventoryId: params.InventoryId,
		CartId:      params.CartId,
	}

	if err := setupQueries(ctx, state); err != nil {
		return err
	}

	inventoryFuture, _ := allocateInventory(ctx, params, state, logger)
	cartFuture, _ := startShoppingCart(ctx, params, state, logger)

	if err := keepShopping(ctx, params, logger); err != nil {
		// the `ContinueAsNew` error could appear here which would early return, hence
		// keeping our "cart" and "inventory" workflows running (since we used ParentClosePolicy.Abandon)
		if !errors.Is(err, workflow.ErrCanceled) {
			return err
		}
	}
	logger.Info("canceling inventory")
	if err := workflow.RequestCancelExternalWorkflow(ctx, params.InventoryId, "").Get(ctx, nil); err != nil {
		logger.Error("failure to cancel inventory", "err", err)
	}
	if err := inventoryFuture.Get(ctx, nil); err != nil {
		logger.Error("inventory cancel failure", "err", err)
	}
	logger.Info("canceling cart")
	if err := workflow.RequestCancelExternalWorkflow(ctx, params.CartId, "").Get(ctx, nil); err != nil {
		logger.Error("failure to cancel cart", "err", err)
	}
	if err := cartFuture.Get(ctx, nil); err != nil {
		logger.Error("cart cancel failure", "err", err)
	}
	logger.Info("session timed out. inventory and cart canceled")
	return nil
}

func allocateInventory(
	ctx workflow.Context,
	params *orchestrations2.StartShopperRequest,
	state *queries.GetShopperResponse,
	logger log.Logger,
) (workflow.ChildWorkflowFuture, workflow.CancelFunc) {
	cctx, cancel := workflow.WithCancel(workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID:          state.InventoryId,
		WaitForCancellation: true,
		// we don't want inventory to go away after a ContinueAsNew
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}))

	f := workflow.ExecuteChildWorkflow(cctx, TypeOrchestrations.Inventory, &orchestrations2.AllocateInventoryRequest{
		Email:       params.Email,
		InventoryId: state.InventoryId,
	})
	return f, cancel
}
func startShoppingCart(
	ctx workflow.Context,
	params *orchestrations2.StartShopperRequest,
	state *queries.GetShopperResponse,
	logger log.Logger,
) (workflow.ChildWorkflowFuture, workflow.CancelFunc) {
	cctx, cancel := workflow.WithCancel(workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID:          params.CartId,
		WaitForCancellation: true,
		// we don't want cart to go away after a ContinueAsNew
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}))

	f := workflow.ExecuteChildWorkflow(cctx, TypeOrchestrations.Cart, &orchestrations2.SetShoppingCartRequest{
		Email:     params.Email,
		CartId:    params.CartId,
		ShopperId: params.ShopperId,
	})
	return f, cancel
}
func keepShopping(ctx workflow.Context, params *orchestrations2.StartShopperRequest, logger log.Logger) error {
	expSess := &ExpiringSession{
		params: params,
	}
	if err := expSess.SleepUntil(
		ctx,
		time.Second*time.Duration(params.DurationSeconds),
		workflow.GetSignalChannel(ctx, SignalName(&commands.RefreshShopperRequest{})),
	); err != nil {
		logger.Error("session errored out", err)
		return err
	}
	return nil
}

func setupQueries(ctx workflow.Context, state *queries.GetShopperResponse) error {
	if err := workflow.SetQueryHandler(
		ctx,
		QueryName(&queries.GetShopperRequest{}),
		func(req *queries.GetShopperRequest) (*queries.GetShopperResponse, error) {
			return state, nil
		},
	); err != nil {
		return fmt.Errorf("failed to setup shopper query %w", err)
	}
	return nil
}
