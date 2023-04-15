package gql

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hashicorp/go-multierror"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/routes"
	"net/http"
)

type Closeable interface {
	Close() error
}
type Handlers struct {
	Handler       http.HandlerFunc
	Playground    http.HandlerFunc
	Subscriptions http.HandlerFunc
	closeables    []Closeable
}

func (h *Handlers) Close() error {
	errors := &multierror.Error{}
	for _, c := range h.closeables {
		if err := c.Close(); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

func NewHandlers(opts ...Option) (*Handlers, error) {

	defaultOpts := []Option{WithPubSub(pubsub.NewPubSub())}
	opts = append(defaultOpts, opts...)
	r := &Resolver{}
	d := graph.DirectiveRoot{}
	c := graph.ComplexityRoot{}

	for _, o := range opts {
		o(r, d, c)
	}
	srv := createGqlHandler(r, d, c)
	subscriptionsSrv := createSubscriptionsHandler(r, d, c)
	playground := playground.Handler("Temporal Shop GraphQL playground", routes.GETGqlPlayground.Raw)
	return &Handlers{
		Handler:       srv,
		Playground:    playground,
		Subscriptions: subscriptionsSrv,
		closeables:    []Closeable{r},
	}, nil
}

func createGqlHandler(r *Resolver, d graph.DirectiveRoot, c graph.ComplexityRoot) http.HandlerFunc {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers:  r,
		Directives: d,
		Complexity: c,
	}))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		logger := log.GetLogger(ctx)
		msg := fmt.Sprintf("gql operation '%s'", oc.OperationName)
		logger.Debug(msg, log.Fields{"q": oc.RawQuery})
		return next(ctx)
	})
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
func createSubscriptionsHandler(r *Resolver, d graph.DirectiveRoot, c graph.ComplexityRoot) http.HandlerFunc {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers:  r,
		Directives: d,
		Complexity: c,
	}))
	srv.AddTransport(transport.SSE{}) // <---- This is the important

	// default server
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
