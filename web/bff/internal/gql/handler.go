package gql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/plugin/federation/testdata/entityresolver/generated"
)

func NewHandlers() (*handler.Server, error) {
	cfg := generated.Config{
		Resolvers:  &Resolver{},
		Directives: generated.DirectiveRoot{},
		Complexity: generated.ComplexityRoot{},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	return srv, nil
}
