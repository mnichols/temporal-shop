package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	temporalClient "github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/app"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/routes"
	"github.com/temporalio/temporal-shop/web/bff/internal/instrumentation/log"
	"logur.dev/logur"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	temporal *temporalClient.Clients
	logger   logur.Logger
	router   *chi.Mux
	cfg      *Config
	inner    *http.Server
}

// NewServer creates a new server with options
// A new root router will be created if one is not provided
func NewServer(ctx context.Context, opts ...Option) (*Server, error) {
	defaultOpts := []Option{
		WithRouter(chi.NewRouter()),
		WithConfig(&Config{}),
	}
	opts = append(defaultOpts, opts...)

	s := &Server{}

	for _, o := range opts {
		o(s)
	}
	s.router.Use(middleware.Logger(s.logger))

	if s.cfg.IsServingUI {
		appHandlers, err := app.NewHandlers(
			app.WithGeneratedAppDirectory(s.cfg.GeneratedAppDir),
			app.WithTemporalClients(s.temporal),
			app.WithMountPath(routes.GETApp.Raw),
		)
		if err != nil {
			return nil, err
		}
		appHandlers.Register(s.router)
	}

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
