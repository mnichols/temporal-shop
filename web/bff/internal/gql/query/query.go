package query

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
)

func NewQuery(t *temporal.Clients) *Query {
	q := &Query{}
	q.temporal = t

	q.shopper = &shopper{t.Client}

	return q
}

type Query struct {
	*shopper
	temporal *temporal.Clients
}
