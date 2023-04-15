package client

import "github.com/temporalio/temporal-shop/services/go/pkg/clients/http"

type Option func(*Client)

func WithHttpClient(h *http.Client) Option {
	return func(c *Client) {
		c.httpClient = h
	}
}
func WithConfig(cfg *Config) Option {
	return func(c *Client) {
		c.cfg = cfg
	}
}
