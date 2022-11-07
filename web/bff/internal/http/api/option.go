package app

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
)

type Option func(h *Handlers)

func WithTemporalClients(c *temporal.Clients) Option {
	return func(h *Handlers) {
		h.temporal = c
	}
}
