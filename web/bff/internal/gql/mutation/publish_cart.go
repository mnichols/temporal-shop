package mutation

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
)

type publishCart struct {
	pubSub *pubsub.PubSub
}

//
//func (m *publishCart) PublishCart(ctx context.Context, input model.PublishCartInput) (*model.PublishCart, error) {
//	logger := log.GetLogger(ctx)
//	logger.Info("PublishCart invoked")
//	out := &model.PublishCart{
//		ID:        input.ID,
//		ShopperID: input.ShopperID,
//		Items:     transformItems(input.Items),
//		Subtotal:  format.Strptr(format.CentsToDollars(int64(input.SubtotalCents))),
//		TaxRate:   format.Strptr(format.BpsToPercentI(input.TaxRateBps)),
//		Total:     format.Strptr(format.CentsToDollarsI(input.TotalCents)),
//		Tax:       format.Strptr(format.CentsToDollarsI(input.TaxCents)),
//	}
//	if err := m.pubSub.PublishCart(ctx, out); err != nil {
//		return nil, err
//	}
//	return out, nil
//}
//
//func transformItems(items []*model.PublishCartItemInput) []*model.CartItem {
//	out := make([]*model.CartItem, len(items))
//	for i := 0; i < len(items); i++ {
//		in := items[i]
//		out[i] = &model.CartItem{
//			ProductID: in.ProductID,
//			Quantity:  in.Quantity,
//			Subtotal:  format.CentsToDollarsI(in.SubtotalCents),
//			Price:     format.CentsToDollarsI(in.PriceCents),
//			Title:     in.Title,
//		}
//	}
//	return out
//}
