package gql

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
)

type Option func(resolver *Resolver, directive graph.DirectiveRoot, complexity graph.ComplexityRoot)

func WithTemporal(c *temporal.Clients) Option {
	return func(resolver *Resolver, _ graph.DirectiveRoot, _ graph.ComplexityRoot) {
		resolver.temporal = c
	}
}
func WithPubSub(p *pubsub.PubSub) Option {
	return func(res *Resolver, _ graph.DirectiveRoot, _ graph.ComplexityRoot) {
		res.pubSub = p
	}
}
func WithTaskQueue(tq string) Option {
	return func(res *Resolver, _ graph.DirectiveRoot, _ graph.ComplexityRoot) {
		res.taskQueue = tq
	}
}
