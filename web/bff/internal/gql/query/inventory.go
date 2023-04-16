package query

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/format"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	sdkclient "go.temporal.io/sdk/client"
	"sort"
)

type inventory struct {
	temporal sdkclient.Client
	shopper  func(ctx context.Context, input *model.ShopperInput) (*model.Shopper, error)
}

func inventoryFromProto(in *queries.GetInventoryResponse) *model.Inventory {
	out := &model.Inventory{
		Categories: []string{""},
	}
	for _, g := range in.Games {
		out.Games = append(out.Games, &model.Game{
			ID:       g.Id,
			Title:    g.Title,
			Category: g.Category,
			ImageURL: g.ImageUrl,
			Price:    format.CentsToDollars(g.PriceCents),
		})
	}
	return out
}
func (q *inventory) Inventory(ctx context.Context, input *model.InventoryInput) (*model.Inventory, error) {
	logger := log.GetLogger(ctx)
	logger.Info("fetching games")
	response := &queries.GetInventoryResponse{}
	s, err := q.shopper(ctx, nil)
	if err != nil {
		return nil, err
	}
	i, err := q.temporal.QueryWorkflow(ctx, s.InventoryID, "", orchestrations.QueryName(&queries.GetInventoryRequest{}))
	if err != nil {
		return nil, err
	}
	if !i.HasValue() {
		return nil, fmt.Errorf("failed to get inventory")
	}
	if ierr := i.Get(response); ierr != nil {
		return nil, fmt.Errorf("failed to get inventory values %w", ierr)
	}
	var games []*model.Game
	categories := make(map[string]interface{})
	for _, g := range response.Games {
		categories[g.Category] = struct{}{}
		if (input == nil || input.Category == nil) || (input != nil && input.Category != nil && *input.Category == g.Category) {
			games = append(games, &model.Game{
				ID:       g.Id,
				Title:    g.Title,
				Category: g.Category,
				ImageURL: g.ImageUrl,
				Price:    format.CentsToDollars(g.PriceCents),
			})
		}
	}
	catSlice := []string{}
	for k := range categories {
		catSlice = append(catSlice, k)
	}
	sort.Strings(catSlice)

	return &model.Inventory{
		Games:      games,
		Categories: catSlice,
	}, nil
}
