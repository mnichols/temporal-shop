package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/auth"
	"logur.dev/logur"
)

type Option func(s *Server)

func WithTemporalClients(c *temporal.Clients) Option {
	return func(s *Server) {
		s.temporal = c
	}
}
func WithRouter(c *chi.Mux) Option {
	return func(s *Server) {
		s.router = c
	}
}

func WithConfig(cfg *Config) Option {
	return func(s *Server) {
		s.cfg = cfg
	}
}
func WithLogger(l logur.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}
func WithAuthenticator(a *auth.Authenticator) Option {
	return func(s *Server) {
		s.authenticator = a
	}
}
