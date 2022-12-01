package orchestrations

import (
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"go.temporal.io/sdk/workflow"
)

func (o *Orchestrations) Shopper(ctx workflow.Context, params *orchestrations2.StartSessionRequest) error {
	return nil
}
