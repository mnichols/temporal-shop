package subscription

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
)

type cart struct {
	pubSub *pubsub.PubSub
}

func (s *cart) Cart(ctx context.Context, input model.CartSubscriptionInput) (<-chan *model.Cart, error) {
	logger := log.GetLogger(ctx)

	auth, ok := middleware.GetAuth(ctx)
	if !ok {
		logger.Error("no auth found")
		return nil, fmt.Errorf("no auth found")
	}
	logger.Info("subscribing!!")
	ch := s.pubSub.SubscribeCart(ctx, auth.SessionID(), input.CartID)
	return ch, nil
}
