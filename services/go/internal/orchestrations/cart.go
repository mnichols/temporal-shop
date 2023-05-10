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
func (w *Orchestrations) Cart(ctx workflow.Context, params *orchestrations.SetShoppingCartRequest) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 2,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	var err error
	calculateRequest := &commands.CalculateShoppingCartRequest{
		CartId:               params.CartId,
		ShopperId:            params.ShopperId,
		TaxRateBps:           shopping.DefaultTaxRateBPS,
		ProductIdsToQuantity: params.ProductIdsToQuantity,
		ProductIdToGame:      map[string]*values.Game{},
	}
	var state *queries.GetCartResponse
	if err = workflow.SetQueryHandler(ctx, QueryName(&queries.GetCartRequest{}), func(req *queries.GetCartRequest) (*queries.GetCartResponse, error) {
		return state, nil
	}); err != nil {
		return fmt.Errorf("failed to setup cart query %w", err)
	}

	if len(params.ProductIdsToQuantity) > 0 {
		var res *inventory.GetGamesResponse
		pids := []string{}
		for pid := range params.ProductIdsToQuantity {
			pids = append(pids, pid)
		}
		if err = workflow.ExecuteActivity(ctx, inventory2.TypeHandlers.GetGames, &inventory.GetGamesRequest{
			Version:           "",
			IncludeProductIds: pids,
		}).Get(ctx, &res); err != nil {
			return err
		}
		for _, g := range res.Games {
			calculateRequest.ProductIdToGame[g.Id] = g
		}
	}
	lao := workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		StartToCloseTimeout: 1 * time.Second,
	})
	if err = workflow.ExecuteLocalActivity(
		lao,
		shopping.TypeHandlers.CalculateShoppingCart,
		calculateRequest,
	).Get(ctx, &state); err != nil {
		return fmt.Errorf("failed to calculate shopping cart %w", err)
	}

	if params.Topic != nil {
		logger.Info("publishing to", "tq", params.Topic.TaskQueue, "activity", params.Topic.Activity)
		publishCtx := workflow.WithTaskQueue(ctx, params.Topic.TaskQueue)
		if pubErr := workflow.ExecuteActivity(publishCtx, params.Topic.Activity, state).Get(publishCtx, nil); pubErr != nil {
			logger.Error("failed to publish", "err", pubErr)
		}
	}

	setCartItemsCommand := &commands.SetCartItemsRequest{}
	setItemsChan := workflow.GetSignalChannel(ctx, SignalName(setCartItemsCommand))
	setItemsChan.Receive(ctx, &setCartItemsCommand)
	logger.Info("set items on cart", "setCartItemsCommand", setCartItemsCommand)

	nextRunParams := &orchestrations.SetShoppingCartRequest{
		CartId:               params.CartId,
		ShopperId:            params.ShopperId,
		Email:                params.Email,
		ProductIdsToQuantity: setCartItemsCommand.ProductIdsToQuantity,
	}
	if setCartItemsCommand.Caller != nil {
		nextRunParams.Topic = &values.Topic{
			TaskQueue: setCartItemsCommand.Caller.TargetTaskQueue,
			Activity:  setCartItemsCommand.Caller.TargetActivity,
		}
	}
	//Drain signal channel asynchronously to avoid signal loss
	for {
		var signalVal string
		ok := setItemsChan.ReceiveAsync(&signalVal)
		if !ok {
			break
		}
		logger.Info("async receipt of signal")
		nextRunParams = &orchestrations.SetShoppingCartRequest{
			CartId:               params.CartId,
			ShopperId:            params.ShopperId,
			Email:                params.Email,
			ProductIdsToQuantity: setCartItemsCommand.ProductIdsToQuantity,
		}
		if setCartItemsCommand.Caller != nil {
			nextRunParams.Topic = &values.Topic{
				TaskQueue: setCartItemsCommand.Caller.TargetTaskQueue,
				Activity:  setCartItemsCommand.Caller.TargetActivity,
			}
		}
	}
	return workflow.NewContinueAsNewError(ctx, TypeOrchestrations.Cart, nextRunParams)
}
