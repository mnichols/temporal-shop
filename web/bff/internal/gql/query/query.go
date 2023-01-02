package query

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
)

func NewQuery(t *temporal.Clients) *Query {
	q := &Query{}
	q.temporal = t

	q.shopper = &shopper{t.Client}
	q.inventory = &inventory{temporal: t.Client, shopper: q.Shopper}

	return q
}

type Query struct {
	*shopper
	*inventory
	temporal *temporal.Clients
}
