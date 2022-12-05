package query

import (
	"context"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
)

type Query struct{}

func (q *Query) Todos(ctx context.Context) ([]*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Query) Dogs(ctx context.Context) ([]*model.Dog, error) {
	//TODO implement me
	panic("implement me")
}
