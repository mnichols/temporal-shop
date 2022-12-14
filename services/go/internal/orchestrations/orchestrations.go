package orchestrations

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/services/go/internal/admin"
	"go.temporal.io/sdk/workflow"
)

var TypeOrchestrations *Orchestrations
var adminHandlers *admin.Handlers

func SignalName(m proto.Message) string {
	if m == nil {
		return ""
	}
	return string(m.ProtoReflect().Descriptor().FullName())
}
func QueryName(m proto.Message) string {
	if m == nil {
		return ""
	}
	return string(m.ProtoReflect().Descriptor().FullName())
}

type Orchestrations struct {
}

func (w *Orchestrations) Ping(ctx workflow.Context, params *orchestrations.PingRequest) (*orchestrations.PingResponse, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var pong *orchestrations.PingResponse
	if err := workflow.ExecuteActivity(ctx, adminHandlers.PingPong, commands.PingRequest{Name: params.Name}).
		Get(ctx, &pong); err != nil {
		return nil, fmt.Errorf("ping pong failed! %w", err)
	}

	return &orchestrations.PingResponse{
		Name: fmt.Sprintf("Pong = %s", pong.Name),
	}, nil
}
