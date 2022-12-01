package auth

import (
	"context"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Sessioner interface {
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}
type TemporalSessionStore struct {
	temporalClient Sessioner
}

func (t *TemporalSessionStore) Validate(ctx context.Context, value string) error {
	refresh := &commands.RefreshShopperRequest{
		LastSeenAt: timestamppb.Now(),
		Email:      value,
	}
	err := t.temporalClient.SignalWorkflow(ctx, value, "", orchestrations.SignalName(refresh), refresh)
	if err != nil {
		return err
	}
	return nil
}
func (t *TemporalSessionStore) Start(ctx context.Context, params *orchestrations2.StartSessionRequest) error {
	logger := log.GetLogger(ctx)
	opts := client.StartWorkflowOptions{
		ID:        params.Id,
		TaskQueue: orchestrations.TaskQueueDefault,
	}
	run, err := t.temporalClient.ExecuteWorkflow(ctx, opts, orchestrations.TypeOrchestrations.Shopper, params)
	if err != nil {
		logger.Error("failed to start session", log.Fields{log.TagError: err})
		return err
	}
	logger.Info("started session", log.Fields{"run": run})
	return nil
}

func NewTemporalSessionStore(temporalClient Sessioner) *TemporalSessionStore {
	return &TemporalSessionStore{
		temporalClient: temporalClient,
	}
}
