package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-multierror"
	temporalClient "github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/api"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/app"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/auth"

	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/routes"
	"github.com/temporalio/temporal-shop/web/bff/internal/instrumentation/log"
	"logur.dev/logur"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	temporal      *temporalClient.Clients
	logger        logur.Logger
	router        *chi.Mux
	cfg           *Config
	inner         *http.Server
	errors        *multierror.Error
	authenticator *auth.Authenticator
}

// NewServer creates a new server with options
// A new root router will be created if one is not provided
func NewServer(ctx context.Context, opts ...Option) (*Server, error) {
	defaultOpts := []Option{
		WithRouter(chi.NewRouter()),
		WithConfig(&Config{}),
	}
	opts = append(defaultOpts, opts...)

	s := &Server{
		errors: &multierror.Error{},
	}

	for _, o := range opts {
		o(s)
	}
	s.router.Use(middleware.Logger(s.logger))

	logger := log.GetLogger(ctx)
	logger.Info("registering routers")

	s.router.Group(s.buildPublicRouter)
	s.router.Group(s.buildSecureRouter)

	s.router.Get("/ping", pingHandler)
	s.router.Get("/health", healthHandler)
	s.router.Get("/readiness", readinessHandler)

	s.inner = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.Port),
		Handler: s.router,
	}
	for _, r := range s.router.Routes() {
		log.GetLogger(ctx).Info("registered the route", log.Fields{"route": r.Pattern})
	}
	return s, nil
}
func (s *Server) appendError(err error) error {
	s.errors = multierror.Append(s.errors, err)
	return s.errors
}

func (s *Server) buildPublicRouter(r chi.Router) {
	var err error
	if s.cfg.IsServingUI {
		var appHandlers *app.Handlers
		if appHandlers, err = app.NewHandlers(
			app.WithGeneratedAppDirectory(s.cfg.GeneratedAppDir),
			app.WithTemporalClients(s.temporal),
			app.WithMountPath(routes.GETApp.Raw),
		); err != nil {
			s.appendError(err)
			return
		}

		appHandlers.Register(r)
	}

}
func (s *Server) buildSecureRouter(r chi.Router) {
	if s.authenticator == nil {
		panic("an authenticator implementation is required for secure router")
	}
	r.Group(func(r chi.Router) {
		r = r.With(middleware.Authenticate(s.authenticator))
		apiHandlers, err := api.NewHandlers(api.WithEncryptionKey(s.cfg.EncryptionKey))
		if err != nil {
			s.appendError(err)
			return
		}
		s.router.Get(routes.GETApi.Raw, apiHandlers.GET)
	})
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	log.GetLogger(ctx).Info("starting http server", log.Fields{
		"port": s.cfg.Port,
	})
	return s.inner.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) {
	if err := s.inner.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown gracefully", logur.Fields{"err": err})
	}
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func pingHandler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		if i, werr := w.Write([]byte(fmt.Sprintf("pong but error %s", err.Error()))); werr != nil {
			fmt.Println("wrote ", i, " bytes", werr)
		}
		return
	}
	if i, werr := w.Write(dump); werr != nil {
		fmt.Println("wrote ", i, "bytes", werr)
	}
}
