package workflows

import (
	"fmt"
	"time"

	"github.com/temporalio/temporal-shop/services/go/internal/admin"
	"github.com/temporalio/temporal-shop/services/go/pkg/messages/commands"
	"github.com/temporalio/temporal-shop/services/go/pkg/messages/events"
	"github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
	"go.temporal.io/sdk/workflow"
)

var TypeWorkflows *Workflows
var adminHandlers *admin.Handlers

type Workflows struct{}

func (w *Workflows) Ping(ctx workflow.Context, params workflows.Ping) (workflows.Pong, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var pong *events.Pong
	if err := workflow.ExecuteActivity(ctx, adminHandlers.PingPong, commands.Ping{
		Value: params.Value,
	}).Get(ctx, &pong); err != nil {
		return workflows.Pong{}, fmt.Errorf("ping pong failed! %w", err)
	}

	return workflows.Pong{
		Value: fmt.Sprintf("Pong = %s", pong.Value),
	}, nil
}
