package gql

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/mutation"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/query"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/subscription"
)

type Resolver struct {
	temporal  *temporal.Clients
	pubSub    *pubsub.PubSub
	taskQueue string
}

func (r Resolver) Mutation() graph.MutationResolver {
	return mutation.NewMutation(
		mutation.WithTemporalClient(r.temporal),
		mutation.WithQuerier(r.Query()),
		mutation.WithPubSub(r.pubSub),
		mutation.WithTaskQueue(r.taskQueue),
	)
}

func (r Resolver) Query() graph.QueryResolver {
	return query.NewQuery(r.temporal)
}
func (r Resolver) Subscription() graph.SubscriptionResolver {
	return subscription.NewSubscription(r.pubSub)
}
func (r Resolver) Close() error {
	if closeable, ok := r.Subscription().(Closeable); ok {
		return closeable.Close()
	}
	return nil
}
