package mutation

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
)

type Mutation struct {
	*setCartItems
	*publishCart
	temporal  *temporal.Clients
	query     graph.QueryResolver
	pubSub    *pubsub.PubSub
	taskQueue string
}

func NewMutation(opts ...Option) *Mutation {
	m := &Mutation{}
	for _, o := range opts {
		o(m)
	}

	m.setCartItems = &setCartItems{m.temporal.Client, m.query, m.taskQueue}
	m.publishCart = &publishCart{
		pubSub: m.pubSub,
	}
	return m
}
