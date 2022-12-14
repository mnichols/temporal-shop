package gql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
)

func NewHandlers(opts ...Option) (*handler.Server, error) {

	r := &Resolver{}
	d := graph.DirectiveRoot{}
	c := graph.ComplexityRoot{}

	for _, o := range opts {
		o(r, d, c)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers:  r,
		Directives: d,
		Complexity: c,
	}))
	return srv, nil
}
