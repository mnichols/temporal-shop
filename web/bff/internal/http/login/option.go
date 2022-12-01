package login

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
)

type Option func(h *Handlers)

func WithTemporalClients(c *temporal.Clients) Option {
	return func(h *Handlers) {
		h.temporal = c
	}
}
func WithEncryptionKey(key string) Option {
	return func(h *Handlers) {
		h.encryptionKey = key
	}
}
func WithSessionStore(a SessionStarter) Option {
	return func(h *Handlers) {
		h.session = a
	}
}
