package temporal

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/internal/admin"
	"github.com/temporalio/temporal-shop/services/go/internal/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/internal/workflows"
	temporalClient "github.com/temporalio/temporal-shop/services/go/pkg/clients/temporal"

	"logur.dev/logur"

	"go.temporal.io/sdk/worker"
)

const NamespaceMessaging = "messaging"
const TaskQueueTemporalShop = "temporal_shop"

type Option func(w *Worker)

func WithTemporal(c *temporalClient.Clients) Option {
	return func(w *Worker) {
		w.temporalClients = c
	}
}

type Config struct{}

type Worker struct {
	temporalClients *temporalClient.Clients
	inner           worker.Worker
	// other clients
}

func NewWorker(_ context.Context, opts ...Option) (*Worker, error) {
	w := &Worker{}
	for _, o := range opts {
		o(w)
	}
	return w, nil
}
func (w *Worker) Shutdown(_ context.Context) {
	// TODO
}
func (w *Worker) register(inner worker.Worker) error {
	wfs := &workflows.Workflows{}
	admin := &admin.Handlers{}
	//repurchasing := &repurchasing2.Handlers{}
	inner.RegisterActivity(admin)
	//inner.RegisterActivity(repurchasing)

	inner.RegisterWorkflow(wfs.Ping)
	//inner.RegisterWorkflow(wfs.EnrollRepurchasingCustomer)
	//inner.RegisterWorkflow(wfs.RemindRepurchasingCustomer)
	return nil
}
func (w *Worker) Start(ctx context.Context) error {
	logger := log.GetLogger(ctx)

	//retention := time.Hour * 24
	//if err := w.temporalClients.NamespaceClient.Register(ctx, &workflowservice.RegisterNamespaceRequest{
	//	Namespace:                        NamespaceMessaging,
	//	WorkflowExecutionRetentionPeriod: &retention,
	//	Description:                      "messaging namespace",
	//}); err != nil {
	//	logger.Error("namespace already registered", logur.Fields{"namespace": NamespaceMessaging, "err": err})
	//}
	inner := worker.New(w.temporalClients.Client, TaskQueueTemporalShop, worker.Options{})

	if err := w.register(inner); err != nil {
		return fmt.Errorf("failed to register workflows/activities: %w", err)
	}

	w.inner = inner
	err := w.inner.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("start worker failed", logur.Fields{"err": err})
		return err
	}
	return nil
}
