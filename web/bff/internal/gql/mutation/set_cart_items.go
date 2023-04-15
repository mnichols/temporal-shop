package mutation

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/subscription"
	sdkclient "go.temporal.io/sdk/client"
)

type setCartItems struct {
	temporal  sdkclient.Client
	q         graph.QueryResolver
	taskQueue string
}

func (a *setCartItems) SetCartItems(ctx context.Context, input *model.SetCartItemsInput) (*model.Cart, error) {
	logger := log.GetLogger(ctx)
	logger.Info("setting items on cart")

	if input == nil {
		return nil, fmt.Errorf("items must be provided")
	}
	if input.CartID == "" {
		shopper, err := a.q.Shopper(ctx, nil)
		if err != nil {
			return nil, err
		}
		input.CartID = shopper.CartID
	}
	cmd := &commands.SetCartItemsRequest{
		CartId:               input.CartID,
		ProductIdsToQuantity: make(map[string]int64),
		Caller: &commands.CallerRequest{
			TargetActivity:  orchestrations.ActivityName(subscription.TypePublishers.PublishCart),
			TargetTaskQueue: a.taskQueue,
		},
	}

	for _, item := range input.Items {
		cmd.ProductIdsToQuantity[item.ProductID] = int64(item.Quantity)
	}
	logger.Info("signaling cart")
	err := a.temporal.SignalWorkflow(ctx, input.CartID, "", orchestrations.SignalName(cmd), cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to set items to cart %w", err)
	}
	return a.q.Cart(ctx, &model.CartInput{
		CartID:    input.CartID,
		ShopperID: "",
	})
}
