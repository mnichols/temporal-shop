package gql

import (
	"github.com/99designs/gqlgen/plugin/federation/testdata/entityresolver/generated"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/mutation"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/query"
)

type Resolver struct {
}

func (r *Resolver) Entity() generated.EntityResolver {
	//TODO implement me
	panic("implement me")
}

func (r Resolver) Mutation() graph.MutationResolver {
	return &mutation.Mutation{}
}

func (r Resolver) Query() graph.QueryResolver {
	return &query.Query{}
}
