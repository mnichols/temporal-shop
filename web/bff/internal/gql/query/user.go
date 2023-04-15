package query

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
)

type user struct {
	q graph.QueryResolver
}

func (u *user) User(ctx context.Context, input *model.UserInput) (*model.User, error) {
	logger := log.GetLogger(ctx)

	auth, ok := middleware.GetAuth(ctx)
	if !ok {
		logger.Error("no auth found")
		return nil, fmt.Errorf("no auth found")
	}
	s, err := u.q.Shopper(ctx, nil)
	if err != nil {
		return nil, err
	}
	token := auth.Token()
	return &model.User{
		Email: s.Email,
		Token: &token,
		Ok:    token != "",
	}, nil
}
