package pubsub

import (
	"context"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	pubsub "github.com/temporalio/temporal-shop/services/go/internal/pubsub/client"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"go.temporal.io/sdk/activity"
	"time"
)

var TypeHandlers *Handlers

type Handlers struct {
	client *pubsub.Client
}

func NewHandlers(client *pubsub.Client) *Handlers {
	return &Handlers{client: client}
}
func (h *Handlers) PublishCart(ctx context.Context, response *queries.GetCartResponse) error {
	info := activity.GetInfo(ctx)
	if time.Now().UTC().Sub(info.ScheduledTime.UTC()) > time.Second*45 {
		// stale publish message so just silently ignore
		return nil
	}
	logger := log.GetLogger(ctx)
	logger = log.WithFields(logger, log.Fields{"cart_id": response.CartId})
	var items []PublishCartItemInput
	for pid, game := range response.ProductIdToGame {
		items = append(items, PublishCartItemInput{
			ProductId:     pid,
			Quantity:      int(response.ProductIdToQuantity[pid]),
			SubtotalCents: int(game.PriceCents * response.ProductIdToQuantity[pid]),
			PriceCents:    int(game.PriceCents),
			Title:         game.Title,
		})
	}
	_, err := PublishCart(ctx, h.client.GraphQLClient(), PublishCartInput{
		Id:            response.CartId,
		ShopperId:     response.ShopperId,
		Items:         items,
		SubtotalCents: int(response.SubtotalCents),
		TaxRateBps:    int(response.TaxRateBps),
		TotalCents:    int(response.TotalCents),
		TaxCents:      int(response.TaxCents),
	})
	if err != nil {
		logger.Error("failed to publish cart")
		return err
	}
	logger.Info("published cart")
	return nil
}
