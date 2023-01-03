package session

import (
	"context"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TemporalSessionStore struct {
	temporalClient Sessioner
}

func (t *TemporalSessionStore) Validate(ctx context.Context, id string) error {
	refresh := &commands.RefreshShopperRequest{
		LastSeenAt: timestamppb.Now(),
		Email:      id,
	}
	err := t.temporalClient.SignalWorkflow(ctx, id, "", orchestrations.SignalName(refresh), refresh)
	if err != nil {
		return err
	}
	return nil
}
func (t *TemporalSessionStore) Start(ctx context.Context, params *orchestrations2.StartShopperRequest) error {
	logger := log.GetLogger(ctx)
	opts := client.StartWorkflowOptions{
		ID:        params.Id,
		TaskQueue: orchestrations.TaskQueueTemporalShop,
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
