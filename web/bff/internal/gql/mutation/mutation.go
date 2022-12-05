package mutation

import (
	"context"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
)

type Mutation struct{}

func (m *Mutation) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}
