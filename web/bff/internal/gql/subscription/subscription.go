package subscription

import "github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"

func NewSubscription(p *pubsub.PubSub) *Subscription {
	result := &Subscription{}
	result.pubSub = p
	result.cart = &cart{
		pubSub: result.pubSub,
	}
	return result
}

type Subscription struct {
	*cart
	pubSub *pubsub.PubSub
}

func (s *Subscription) Close() error {
	return s.pubSub.Close()
}
