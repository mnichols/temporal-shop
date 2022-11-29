package auth

import (
	"context"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionSignaler interface {
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
}
type TemporalSessionStore struct {
	temporalClient SessionSignaler
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

func NewTemporalSessionStore(temporalClient SessionSignaler) *TemporalSessionStore {
	return &TemporalSessionStore{
		temporalClient: temporalClient,
	}
}
