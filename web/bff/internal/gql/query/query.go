package query

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
)

func NewQuery(t *temporal.Clients) *Query {
	return &Query{temporal: t}
}

type Query struct {
	temporal *temporal.Clients
}

func (q *Query) Games(ctx context.Context) (*model.Games, error) {
	logger := log.GetLogger(ctx)
	logger.Info("fetching games")
	auth, ok := middleware.GetAuth(ctx)
	if !ok {
		logger.Error("no auth found")
		return nil, fmt.Errorf("no auth found")
	}
	shopper := &queries.GetShopperResponse{}
	inventory := &queries.GetInventoryResponse{}

	me, err := q.temporal.Client.QueryWorkflow(ctx, auth.SessionID.String(), "", orchestrations.QueryName(&queries.GetShopperRequest{}))
	if err != nil {
		return nil, err
	}
	if !me.HasValue() {
		return nil, fmt.Errorf("failed to get shopper")
	}
	if merr := me.Get(shopper); merr != nil {
		return nil, fmt.Errorf("failed to get shopper values %w", merr)
	}
	logger.Info("hydrated shopper", log.Fields{"shopper": shopper.Email})

	i, err := q.temporal.Client.QueryWorkflow(ctx, shopper.InventoryId, "", orchestrations.QueryName(&queries.GetInventoryRequest{}))
	if err != nil {
		return nil, err
	}
	if !i.HasValue() {
		return nil, fmt.Errorf("failed to get inventory")
	}
	if ierr := i.Get(inventory); ierr != nil {
		return nil, fmt.Errorf("failed to get inventory values %w", ierr)
	}

	logger.Info("hydrated inventory", log.Fields{"inventory": inventory.InventoryId})

	games := []*model.Game{}
	for _, g := range inventory.Games {
		games = append(games, &model.Game{
			ID:       g.Id,
			Product:  g.Product,
			Category: g.Category,
			ImageURL: g.ImageUrl,
			Price:    g.Price,
		})
	}

	logger.Info("got the games!")
	return &model.Games{
		Items: games,
	}, nil
}

func (q *Query) Shopper(ctx context.Context) (*model.Shopper, error) {
	//TODO implement me
	panic("implement me")
}
