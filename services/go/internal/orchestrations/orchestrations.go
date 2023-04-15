package orchestrations

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
	"runtime"
	"strings"
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
func ActivityName(activity interface{}) string {
	name, _ := getFunctionName(activity)
	return name
}

// lifted from sdk-go
func getFunctionName(i interface{}) (name string, isMethod bool) {
	if fullName, ok := i.(string); ok {
		return fullName, false
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	// Full function name that has a struct pointer receiver has the following format
	// <prefix>.(*<type>).<function>
	isMethod = strings.ContainsAny(fullName, "*")
	elements := strings.Split(fullName, ".")
	shortName := elements[len(elements)-1]
	// This allows to call activities by method pointer
	// Compiler adds -fm suffix to a function name which has a receiver
	// Note that this works even if struct pointer used to get the function is nil
	// It is possible because nil receivers are allowed.
	// For example:
	// var a *Activities
	// ExecuteActivity(ctx, a.Foo)
	// will call this function which is going to return "Foo"
	return strings.TrimSuffix(shortName, "-fm"), isMethod
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
