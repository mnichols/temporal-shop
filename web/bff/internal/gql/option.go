package gql

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
)

type Option func(resolver *Resolver, directive graph.DirectiveRoot, complexity graph.ComplexityRoot)

func WithTemporal(c *temporal.Clients) Option {
	return func(resolver *Resolver, _ graph.DirectiveRoot, _ graph.ComplexityRoot) {
		resolver.temporal = c
	}
}
