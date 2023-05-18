package admin

import (
	"context"
	commands "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/commands/v1"
)

type Handlers struct{}

func (h *Handlers) PingPong(ctx context.Context, cmd *commands.PingRequest) (*commands.PingResponse, error) {
	return &commands.PingResponse{
		Name: cmd.Name,
	}, nil
}
