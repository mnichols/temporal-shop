package orchestrations

import (
	"github.com/temporalio/temporal-shop/services/go/internal/orchestrations"
	"go.temporal.io/api/enums/v1"
	"google.golang.org/protobuf/proto"
)

var WorkflowStatusRunning = enums.WORKFLOW_EXECUTION_STATUS_RUNNING
var TypeOrchestrations = orchestrations.TypeOrchestrations
var TaskQueueDefault = orchestrations.TaskQueueDefault

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
