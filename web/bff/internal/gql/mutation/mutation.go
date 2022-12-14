package mutation

import (
	"context"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
)

type Mutation struct{}

func (m Mutation) AddGameToCart(ctx context.Context, input model.CartItem) (*model.Cart, error) {
	//TODO implement me
	panic("implement me")
}
