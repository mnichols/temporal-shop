package clients

import (
	"github.com/hashicorp/go-multierror"
	pubsub "github.com/temporalio/temporal-shop/services/go/internal/pubsub/client"
	"github.com/temporalio/temporal-shop/services/go/pkg/clients/http"
	inventory2 "github.com/temporalio/temporal-shop/services/go/pkg/clients/inventory"
	"github.com/temporalio/temporal-shop/services/go/pkg/clients/temporal"
	"logur.dev/logur"
)

type Option func(*Clients)

func WithTemporal(t *temporal.Clients, err error) Option {
	return func(c *Clients) {
		c.temporal = t
		c.clientErrors = multierror.Append(c.clientErrors, err)
	}
}
func WithInventory(i *inventory2.Client, err error) Option {
	return func(c *Clients) {
		c.inventory = i
		c.clientErrors = multierror.Append(c.clientErrors, err)
	}
}
func WithHttp(h *http.Client, err error) Option {
	return func(c *Clients) {
		c.http = h
		c.clientErrors = multierror.Append(c.clientErrors, err)
	}
}
func WithPubSub(p *pubsub.Client, err error) Option {
	return func(c *Clients) {
		c.pubSub = p
		c.clientErrors = multierror.Append(c.clientErrors, err)
	}
}
func WithLogger(l logur.Logger) Option {
	return func(c *Clients) {
		c.logger = l
	}
}
