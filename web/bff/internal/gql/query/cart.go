package query

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/format"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	sdkclient "go.temporal.io/sdk/client"
	"sort"
	"time"
)

type cart struct {
	temporal sdkclient.Client
	q        graph.QueryResolver
}

func (a *cart) Cart(ctx context.Context, input *model.CartInput) (*model.Cart, error) {
	logger := log.GetLogger(ctx)

	if input == nil || input.CartID == "" {
		s, err := a.q.Shopper(ctx, nil)
		if err != nil {
			return nil, err
		}
		input = &model.CartInput{ShopperID: s.ID, CartID: s.CartID}
	}
	query := &queries.GetCartRequest{
		CartId: input.CartID,
	}
	response := &queries.GetCartResponse{}
	logger.Debug("requesting cart details", log.Fields{"cart_id": input.CartID})
	val, err := a.temporal.QueryWorkflow(ctx, input.CartID, "", orchestrations.QueryName(query), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query val %w", err)
	}
	if !val.HasValue() {
		return nil, fmt.Errorf("failed to get response")
	}
	if merr := val.Get(response); merr != nil {
		return nil, fmt.Errorf("failed to get response values %w", merr)
	}
	ts := time.Now().UTC()
	// TODO format nums as strings
	result := &model.Cart{
		ID:        response.CartId,
		ShopperID: response.ShopperId,
		Items:     []*model.CartItem{},
		Total:     format.Strptr(format.CentsToDollars(response.TotalCents)),
		Subtotal:  format.Strptr(format.CentsToDollars(response.SubtotalCents)),
		TaxRate:   format.Strptr(format.BpsToPercent(response.TaxRateBps)),
		Tax:       format.Strptr(format.CentsToDollars(response.TotalCents - response.SubtotalCents)),
		Timestamp: &ts,
	}
	for pid, qty := range response.ProductIdToQuantity {
		g := response.ProductIdToGame[pid]
		result.Items = append(result.Items, &model.CartItem{
			ProductID: pid,
			Quantity:  int(qty),
			Subtotal:  format.CentsToDollars(g.PriceCents * qty),
			Price:     format.CentsToDollars(g.PriceCents),
			Title:     g.Title,
			Category:  g.Category,
		})
	}
	sort.Slice(result.Items, func(i, j int) bool {
		x := result.Items[i]
		y := result.Items[j]
		sortByCategory := x.Category < y.Category
		if x.Category == y.Category {
			return x.Title < y.Title
		}
		return sortByCategory
	})
	return result, nil
}
