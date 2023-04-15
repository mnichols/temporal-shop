package orchestrations

import (
	"github.com/temporalio/temporal-shop/services/go/internal/orchestrations"
	temporal_shop "github.com/temporalio/temporal-shop/services/go/internal/workers/temporal"
	"go.temporal.io/api/enums/v1"
)

var WorkflowStatusRunning = enums.WORKFLOW_EXECUTION_STATUS_RUNNING
var TypeOrchestrations = orchestrations.TypeOrchestrations
var TaskQueueTemporalShop = temporal_shop.TaskQueueTemporalShop
var SignalName = orchestrations.SignalName
var QueryName = orchestrations.QueryName
var ActivityName = orchestrations.ActivityName
