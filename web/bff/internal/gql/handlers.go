package gql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/routes"
	"net/http"
)

type Handlers struct {
	Handler    http.HandlerFunc
	Playground http.HandlerFunc
}

func NewHandlers(opts ...Option) (*Handlers, error) {

	r := &Resolver{}
	d := graph.DirectiveRoot{}
	c := graph.ComplexityRoot{}

	for _, o := range opts {
		o(r, d, c)
	}

	srv := createGqlHandler(r, d, c)
	playground := playground.Handler("Temporal Shop GraphQL playground", routes.GETGqlPlayground.Raw)

	return &Handlers{
		Handler:    srv,
		Playground: playground,
	}, nil
}

func createGqlHandler(r *Resolver, d graph.DirectiveRoot, c graph.ComplexityRoot) http.HandlerFunc {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers:  r,
		Directives: d,
		Complexity: c,
	}))
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
