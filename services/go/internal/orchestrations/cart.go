package orchestrations

import (
	"fmt"
	"github.com/temporalio/temporal-shop/api/inventory/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/values/v1"
	inventory2 "github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"github.com/temporalio/temporal-shop/services/go/internal/shopping"
	"go.temporal.io/sdk/workflow"
	"time"
)

// Cart is an entity workflow for a Shopping Cart
// TODO guard against signal flood, doing continueAsNew after N signals
func (w *Orchestrations) Cart(ctx workflow.Context, params *orchestrations.StartShoppingCartRequest) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 2,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	cart := shopping.NewShoppingCart(shopping.ShoppingCartArgs{
		ID:         params.CartId,
		ShopperID:  params.ShopperId,
		TaxRateBPS: shopping.DefaultTaxRateBPS,
	})

	var err error
	var state *queries.GetCartResponse
	state, err = cart.Empty()
	if err != nil {
		return err
	}
	// todo decide about payment sealing cart
	paymentStarted := false
	if err = workflow.SetQueryHandler(ctx, QueryName(&queries.GetCartRequest{}), func(req *queries.GetCartRequest) (*queries.GetCartResponse, error) {
		return state, nil
	}); err != nil {
		return fmt.Errorf("failed to setup cart query %w", err)
	}

	addItemsChan := workflow.GetSignalChannel(ctx, SignalName(&commands.SetCartItemsRequest{}))
	cartOpsCtx, cancelOps := workflow.WithCancel(ctx)
	workflow.Go(cartOpsCtx, func(ctx workflow.Context) {
		// this implementation must take care to use the appropriate context, here the `cartOpsCtx`
		// or else Temporal will fail with `Trying to block on a coroutine which is already blocked` error.
		fetchGames := func(pids []string) ([]*values.Game, error) {
			var res *inventory.GetGamesResponse
			if err := workflow.ExecuteActivity(ctx, inventory2.TypeHandlers.GetGames, &inventory.GetGamesRequest{
				Version:           "",
				IncludeProductIds: pids,
			}).Get(ctx, &res); err != nil {
				return nil, err
			}
			return res.Games, nil
		}
		logger := workflow.GetLogger(ctx)
		for {
			var callerRequest *commands.CallerRequest
			sel := workflow.NewSelector(ctx)
			sel.AddReceive(addItemsChan, func(c workflow.ReceiveChannel, more bool) {
				cmd := &commands.SetCartItemsRequest{}
				c.Receive(ctx, &cmd)
				logger.Info("set items on cart", "cmd", cmd)
				cart.Append(cmd)
				callerRequest = cmd.Caller
			})
			sel.Select(ctx)
			state, err = cart.Calculate(fetchGames)
			// TODO tune retry options on this to be two shots only
			if callerRequest != nil {
				logger.Info("publishing to", "tq", callerRequest.TargetTaskQueue)
				publishCtx := workflow.WithTaskQueue(ctx, callerRequest.TargetTaskQueue)
				if pubErr := workflow.ExecuteActivity(publishCtx, callerRequest.TargetActivity, state).Get(publishCtx, nil); pubErr != nil {
					logger.Error("failed to publish", "err", pubErr)
				}
			}
			if err != nil {
				logger.Error("state calculation failed. canceling ops", "err", err)
				// todo capture error
				cancelOps()
				break
			}
		}
	})
	// stayin alive
	workflow.Await(ctx, func() bool {
		return paymentStarted || ctx.Err() != nil
	})
	cancelOps()
	return nil
}
