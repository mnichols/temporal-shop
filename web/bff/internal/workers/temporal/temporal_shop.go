package temporal

import (
	"context"
	"fmt"
	temporalClient "github.com/temporalio/temporal-shop/services/go/pkg/clients/temporal"
	pubsub "github.com/temporalio/temporal-shop/web/bff/internal/gql/pubsub"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/subscription"

	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"

	"logur.dev/logur"

	"go.temporal.io/sdk/worker"
)

const TaskQueueTemporalShop = "temporal_shop"

type Option func(w *Worker)

func WithConfig(cfg *Config) Option {
	return func(w *Worker) {
		w.cfg = cfg
	}
}
func WithTemporal(c *temporalClient.Clients) Option {
	return func(w *Worker) {
		w.temporalClients = c
	}
}

func WithPubSub(p *pubsub.PubSub) Option {
	return func(w *Worker) {
		w.pubSub = p
	}
}

type Worker struct {
	temporalClients *temporalClient.Clients
	inner           worker.Worker
	pubSub          *pubsub.PubSub
	cfg             *Config

	// other clients
}

func NewWorker(_ context.Context, opts ...Option) (*Worker, error) {
	w := &Worker{}
	defaultOpts := []Option{}
	opts = append(defaultOpts, opts...)
	for _, o := range opts {
		o(w)
	}
	return w, nil
}
func (w *Worker) Shutdown(_ context.Context) {
	// TODO
}
func (w *Worker) register(inner worker.Worker) error {
	pbsb := subscription.NewPublishers(w.pubSub)

	inner.RegisterActivity(pbsb)
	return nil
}
func (w *Worker) Start(ctx context.Context) error {
	logger := log.WithFields(log.GetLogger(ctx), log.Fields{"task_queue": w.cfg.TaskQueue})

	inner := worker.New(w.temporalClients.Client, w.cfg.TaskQueue, worker.Options{})

	if err := w.register(inner); err != nil {
		return fmt.Errorf("failed to register workflows/activities: %w", err)
	}

	w.inner = inner
	logger.Info("starting worker")
	err := w.inner.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("start worker failed", logur.Fields{"err": err})
		return err
	}
	return nil
}
