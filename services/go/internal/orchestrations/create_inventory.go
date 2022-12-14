package orchestrations

import (
	"fmt"
	"github.com/temporalio/temporal-shop/api/inventory/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	inventory2 "github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"go.temporal.io/sdk/workflow"
	"time"
)

func (w *Orchestrations) CreateInventory(ctx workflow.Context, params *orchestrations.CreateInventoryRequest) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 2,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	state := &queries.GetInventoryResponse{InventoryId: params.Id, Games: []*queries.Game{}}
	if err := workflow.SetQueryHandler(ctx, QueryName(&queries.GetInventoryRequest{}), func(req *queries.GetInventoryRequest) (*queries.GetInventoryResponse, error) {
		return state, nil
	}); err != nil {
		return fmt.Errorf("failed to setup shopper query %w", err)
	}
	var games *inventory.GetGamesResponse
	if err := workflow.ExecuteActivity(ctx, inventory2.TypeHandlers.GetGames, &inventory.GetGamesRequest{Version: "1"}).Get(ctx, &games); err != nil {
		return fmt.Errorf("failed to get games %w", err)
	}
	for _, g := range games.Games {
		state.Games = append(state.Games, &queries.Game{
			Id:       g.Id,
			Product:  g.Product,
			Category: g.Category,
			ImageUrl: g.ImageUrl,
			Price:    g.Price,
		})
	}
	// stayin alive
	workflow.Await(ctx, func() bool {
		return ctx.Err() != nil
	})
	return nil
}
