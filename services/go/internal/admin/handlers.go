package admin

import (
	"context"
)

type Handlers struct{}

func (h *Handlers) PingPong(ctx context.Context, cmd *commands.PingRequest) (*commands.PingResponse, error) {
	return &commands.PingResponse{
		Name: cmd.Name,
	}, nil
}
