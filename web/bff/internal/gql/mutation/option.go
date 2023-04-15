package mutation

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
)

type Option func(mutation *Mutation)

func WithTemporalClient(t *temporal.Clients) Option {
	return func(m *Mutation) {
		m.temporal = t
	}
}
func WithQuerier(r graph.QueryResolver) Option {
	return func(m *Mutation) {
		m.query = r
	}
}
func WithPubSub(p *pubsub.PubSub) Option {
	return func(m *Mutation) {
		m.pubSub = p
	}
}
func WithTaskQueue(tq string) Option {
	return func(m *Mutation) {
		m.taskQueue = tq
	}
}
