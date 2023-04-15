package query

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
	sdkclient "go.temporal.io/sdk/client"
)

type shopper struct {
	temporal sdkclient.Client
}

func (q *shopper) Shopper(ctx context.Context, input *model.ShopperInput) (*model.Shopper, error) {
	logger := log.GetLogger(ctx)
	auth, ok := middleware.GetAuth(ctx)
	if !ok {
		logger.Error("no auth found")
		return nil, fmt.Errorf("no auth found")
	}

	response := &queries.GetShopperResponse{}

	me, err := q.temporal.QueryWorkflow(ctx, auth.SessionID(), "", orchestrations.QueryName(&queries.GetShopperRequest{}))
	if err != nil {
		return nil, err
	}
	if !me.HasValue() {
		return nil, fmt.Errorf("failed to get response")
	}
	if merr := me.Get(response); merr != nil {
		return nil, fmt.Errorf("failed to get response values %w", merr)
	}
	return &model.Shopper{
		ID:          response.ShopperId,
		Email:       response.Email,
		InventoryID: response.InventoryId,
		CartID:      response.CartId,
	}, nil
}
