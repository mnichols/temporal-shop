package temporal_shop

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/api/inventory/v1"
	"github.com/temporalio/temporal-shop/services/go/internal/admin"
	inventory2 "github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"github.com/temporalio/temporal-shop/services/go/internal/orchestrations"
	pubsubClient "github.com/temporalio/temporal-shop/services/go/internal/pubsub/client"

	invClient "github.com/temporalio/temporal-shop/services/go/pkg/clients/inventory"
	temporalClient "github.com/temporalio/temporal-shop/services/go/pkg/clients/temporal"

	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"

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
func WithInventoryClient(c *invClient.Client) Option {
	return func(w *Worker) {
		w.inventoryClient = c
	}
}
func WithPubSub(p *pubsubClient.Client) Option {
	return func(w *Worker) {
		w.pubSub = p
	}
}

type Config struct{}

type Worker struct {
	temporalClients *temporalClient.Clients
	inner           worker.Worker
	inventoryClient inventory.InventoryServiceClient
	pubSub          *pubsubClient.Client

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
	wfs := &orchestrations.Orchestrations{}
	admin := &admin.Handlers{}
	inventory := inventory2.NewHandlers(w.inventoryClient)
	inner.RegisterActivity(inventory)
	inner.RegisterActivity(admin)

	inner.RegisterWorkflow(wfs.Ping)
	inner.RegisterWorkflow(wfs.Shopper)
	inner.RegisterWorkflow(wfs.Inventory)
	inner.RegisterWorkflow(wfs.Cart)
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
	inner := worker.New(w.temporalClients.Client, TaskQueueTemporalShop, worker.Options{
		MaxConcurrentWorkflowTaskPollers: 4,
	})

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
