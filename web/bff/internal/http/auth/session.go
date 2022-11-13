package auth

import (
	"context"
	"go.temporal.io/api/workflowservice/v1"
)

type WorkflowExecutionDescriber interface {
	DescribeWorkflowExecution(context.Context, string, string) (*workflowservice.DescribeWorkflowExecutionResponse, error)
}
type TemporalSessionStore struct {
	temporalClient WorkflowExecutionDescriber
}

func (t *TemporalSessionStore) Validate(ctx context.Context, value string) error {

	_, err := t.temporalClient.DescribeWorkflowExecution(ctx, value, "")
	if err != nil {
		return err
	}
	return nil
}

func NewTemporalSessionStore(temporalClient WorkflowExecutionDescriber) *TemporalSessionStore {
	return &TemporalSessionStore{
		temporalClient: temporalClient,
	}
}
