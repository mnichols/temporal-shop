package temporal

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
)

type Config struct {
	TaskQueue string
}

func (c *Config) Prefix() string {
	return "temporalshop"
}
func (c *Config) Override() error {
	if c.TaskQueue == "" {
		c.TaskQueue = temporal.GetIdentity("")
	}
	return nil
}
