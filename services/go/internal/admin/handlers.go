package admin

import (
	"context"

	"github.com/temporalio/temporal-shop/services/go/pkg/messages/commands"
	"github.com/temporalio/temporal-shop/services/go/pkg/messages/events"
)

type Handlers struct{}

func (h *Handlers) PingPong(ctx context.Context, cmd commands.Ping) (events.Pong, error) {
	return events.Pong{
		Value: cmd.Value,
	}, nil
}
