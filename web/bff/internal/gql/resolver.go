package gql

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/mutation"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/query"
)

type Resolver struct {
	temporal *temporal.Clients
}

func (r Resolver) Mutation() graph.MutationResolver {
	return &mutation.Mutation{}
}

func (r Resolver) Query() graph.QueryResolver {
	return query.NewQuery(r.temporal)
}
