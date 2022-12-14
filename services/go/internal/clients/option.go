package clients

import (
	"github.com/hashicorp/go-multierror"
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

func WithLogger(l logur.Logger) Option {
	return func(c *Clients) {
		c.logger = l
	}
}
