package auth

import (
	"context"
	sdkclient "go.temporal.io/sdk/client"
)

type TemporalSessionStore struct {
	temporalClient sdkclient.Client
}

func (t *TemporalSessionStore) Validate(ctx context.Context, value string) error {
	_, err := t.temporalClient.DescribeWorkflowExecution(ctx, value, "")
	if err != nil {
		return err
	}
	return nil
}

func NewTemporalSessionStore(c sdkclient.Client) *TemporalSessionStore {
	return &TemporalSessionStore{
		temporalClient: c,
	}
}
