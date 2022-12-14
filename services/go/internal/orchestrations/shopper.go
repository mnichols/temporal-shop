package orchestrations

import (
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
	"time"
)

const SessionMinSeconds = 3600 * 24 * 60

type ShopperState struct {
	InventoryID string
	ShopperID   string
	Email       string
}

// ExpiringSession is an example of a timer that can have its wake time updated
type ExpiringSession struct {
	durationSeconds time.Duration
}

// SleepUntil sleeps until the provided wake-up time.
// The wake-up time can be updated at any time by sending a new time over updateWakeUpTimeCh.
// Supports ctx cancellation.
// Returns temporal.CanceledError if ctx was canceled.
func (u *ExpiringSession) SleepUntil(ctx workflow.Context, durationSeconds time.Duration, refreshShopperChan workflow.ReceiveChannel) (err error) {
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
				u.durationSeconds = time.Second * time.Duration(req.DurationSeconds)
				logger.Info("Wake up time update requested")
			}).
			Select(timerCtx)
	}
	return ctx.Err()
}

func (w *Orchestrations) Shopper(ctx workflow.Context, params *orchestrations2.StartSessionRequest) error {
	if params.DurationSeconds == 0 {
		params.DurationSeconds = SessionMinSeconds
	}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 3,
	})
	logger := log.With(workflow.GetLogger(ctx), "email", params.Email)
	logger.Info("started shopper")
	state := &queries.GetShopperResponse{ShopperId: params.Id, Email: params.Email, InventoryId: fmt.Sprintf("inv_%s", params.Id)}
	if err := workflow.SetQueryHandler(ctx, QueryName(&queries.GetShopperRequest{}), func(req *queries.GetShopperRequest) (*queries.GetShopperResponse, error) {
		return state, nil
	}); err != nil {
		return fmt.Errorf("failed to setup shopper query %w", err)
	}
	cctx, cancelInventory := workflow.WithCancel(workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: state.InventoryId,
	}))
	expSess := &ExpiringSession{}

	workflow.ExecuteChildWorkflow(cctx, TypeOrchestrations.CreateInventory, &orchestrations2.CreateInventoryRequest{Email: params.Email, Id: state.InventoryId})
	logger.Debug("inventory created", "inv")
	if err := expSess.SleepUntil(
		ctx,
		time.Second*time.Duration(params.DurationSeconds),
		workflow.GetSignalChannel(ctx, SignalName(&commands.RefreshShopperRequest{})),
	); err != nil {
		logger.Error("session errored out", err)
		return err
	}
	cancelInventory()
	logger.Info("session timed out and inventory canceled")
	return nil
}
