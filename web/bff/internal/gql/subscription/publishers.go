package subscription

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/format"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
	"time"
)

var TypePublishers *Publishers

type Publishers struct {
	pubSub *pubsub.PubSub
}

func NewPublishers(p *pubsub.PubSub) *Publishers {
	return &Publishers{p}
}

func (p *Publishers) PublishCart(ctx context.Context, cmd *queries.GetCartResponse) error {
	logger := log.GetLogger(ctx)
	logger.Info("PublishCart invoked")
	items, err := transformItems(cmd)
	if err != nil {
		return err
	}
	ts := time.Now().UTC()
	out := &model.Cart{
		ID:        cmd.CartId,
		ShopperID: cmd.ShopperId,
		Items:     items,
		Subtotal:  format.Strptr(format.CentsToDollars(cmd.SubtotalCents)),
		TaxRate:   format.Strptr(format.BpsToPercentI(int(cmd.TaxRateBps))),
		Total:     format.Strptr(format.CentsToDollarsI(int(cmd.TotalCents))),
		Tax:       format.Strptr(format.CentsToDollarsI(int(cmd.TaxCents))),
		Timestamp: &ts,
	}
	return p.pubSub.PublishCart(ctx, out)
}
func transformItems(cmd *queries.GetCartResponse) ([]*model.CartItem, error) {
	out := make([]*model.CartItem, len(cmd.ProductIdToQuantity))
	i := 0
	for pid, qty := range cmd.ProductIdToQuantity {
		game, exists := cmd.ProductIdToGame[pid]
		if !exists {
			return nil, fmt.Errorf("could not find game %v", pid)
		}
		out[i] = &model.CartItem{
			ProductID: game.Id,
			Quantity:  int(qty),
			Subtotal:  format.CentsToDollarsI(int(game.PriceCents * qty)),
			Price:     format.CentsToDollarsI(int(game.PriceCents)),
			Title:     game.Title,
		}
		i = i + 1
	}
	return out, nil
}
