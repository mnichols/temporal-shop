package orchestrations

import (
	"fmt"
	inventory "github.com/temporalio/temporal-shop/services/go/api/generated/inventory/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/orchestrations/v1"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	values "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/values/v1"
	inventory2 "github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"go.temporal.io/sdk/workflow"
	"time"
)

func (w *Orchestrations) Inventory(ctx workflow.Context, params *orchestrations2.AllocateInventoryRequest) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 2,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	state := &queries.GetInventoryResponse{InventoryId: params.InventoryId, Games: []*values.Game{}}
	if err := workflow.SetQueryHandler(ctx, QueryName(&queries.GetInventoryRequest{}), func(req *queries.GetInventoryRequest) (*queries.GetInventoryResponse, error) {
		return state, nil
	}); err != nil {
		return fmt.Errorf("failed to setup inventory query %w", err)
	}
	workflow.Go(ctx, func(ctx) {
		if err := workflow.ExecuteActivity(ctx, inventory2.TypeHandlers.GetGames, &inventory.GetGamesRequest{Version: "1"}).Get(ctx, &games); err != nil {
			return fmt.Errorf("failed to get games %w", err)
		}
	})

	var games *inventory.GetGamesResponse
	if err := workflow.ExecuteActivity(ctx, inventory2.TypeHandlers.GetGames, &inventory.GetGamesRequest{Version: "1"}).Get(ctx, &games); err != nil {
		return fmt.Errorf("failed to get games %w", err)
	}
	for _, g := range games.Games {
		state.Games = append(state.Games, g)
	}
	// stayin alive
	if err := workflow.Await(ctx, func() bool {
		return ctx.Err() != nil
	}); err != nil {
		fmt.Printf("inventory was canceled %v", err)
		return err
	}
	return nil
}
