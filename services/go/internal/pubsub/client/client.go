package client

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"github.com/temporalio/temporal-shop/services/go/pkg/clients/http"
)

type Client struct {
	inner      graphql.Client
	httpClient *http.Client
	cfg        *Config
}

func (c *Client) GraphQLClient() graphql.Client {
	return c.inner
}

func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	c := &Client{}
	for _, o := range opts {
		o(c)
	}
	inner := graphql.NewClient(c.cfg.HostPort, c.httpClient)
	c.inner = inner
	return c, nil
}
