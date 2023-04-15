package query

import (
	"context"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
)

type ping struct {
}

func (p *ping) Ping(ctx context.Context, input *model.PingInput) (*model.Pong, error) {
	return &model.Pong{
		Timestamp: input.Timestamp,
		Value:     input.Value,
	}, nil
}
