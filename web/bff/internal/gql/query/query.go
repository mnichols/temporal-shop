package query

import (
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
)

func NewQuery(t *temporal.Clients) *Query {
	q := &Query{}
	q.temporal = t

	q.shopper = &shopper{t.Client}
	q.inventory = &inventory{temporal: t.Client, shopper: q.Shopper}
	q.cart = &cart{temporal: t.Client, q: q}
	q.user = &user{q: q}
	q.ping = &ping{}
	return q
}

type Query struct {
	*shopper
	*inventory
	*cart
	*user
	*ping
	temporal *temporal.Clients
}
