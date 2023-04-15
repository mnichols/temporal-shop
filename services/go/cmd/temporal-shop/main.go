package main

import (
	"context"
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/internal/grpc"
	"github.com/temporalio/temporal-shop/services/go/internal/inventory"
	pubsub "github.com/temporalio/temporal-shop/services/go/internal/pubsub/client"
	"github.com/temporalio/temporal-shop/services/go/pkg/clients/http"
	inventory2 "github.com/temporalio/temporal-shop/services/go/pkg/clients/inventory"
	temporal2 "github.com/temporalio/temporal-shop/services/go/pkg/clients/temporal"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/temporalio/temporal-shop/services/go/internal/clients"
	"github.com/temporalio/temporal-shop/services/go/internal/workers/temporal"
	"github.com/temporalio/temporal-shop/services/go/pkg/config"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"golang.org/x/sync/errgroup"
	"logur.dev/logur"
)

type startable interface {
	Start(context.Context) error
	Shutdown(context.Context)
}

type appConfig struct {
	Log            *log.Config
	TemporalClient *temporal2.Config
	TemporalWorker *temporal_shop.Config
	HTTPConfig     *http.Config
	GRPCConfig     *grpc.Config
	PubSub         *pubsub.Config
}

func main() {
	// config root
	config.MustLoad()
	var err error

	appCfg := &appConfig{}
	config.MustUnmarshalAll(appCfg)

	ctx, done := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// set up signal listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	// logging
	var logger logur.Logger
	if logger, err = log.NewLogger(ctx, appCfg.Log); err != nil {
		panic("failed to create logger" + err.Error())
	}
	ctx = log.WithLogger(ctx, logger)
	// apps
	clientConn, err := grpc.NewClientConnection(ctx, appCfg.GRPCConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create grpc client conn %s", err.Error()))
	}

	pubsubHTTPClient, err := http.NewClient(ctx, http.WithPlugin(http.NewBasicAuth(appCfg.PubSub.Username, appCfg.PubSub.Password)))
	if err != nil {
		panic(fmt.Sprintf("failed to create http client %s", err.Error()))
	}
	// clients
	clients := clients.MustGetClients(ctx,
		clients.WithPubSub(pubsub.NewClient(ctx, pubsub.WithHttpClient(pubsubHTTPClient), pubsub.WithConfig(appCfg.PubSub))),
		clients.WithTemporal(temporal2.NewClients(ctx,
			temporal2.WithConfig(appCfg.TemporalClient),
			temporal2.WithLogger(logger))),
		clients.WithInventory(inventory2.NewClient(ctx, clientConn)),
	)

	defer func() {
		if perr := clients.Close(); perr != nil {
			logger.Error("failed to close clients gracefully", logur.Fields{"err": perr})
		}
	}()

	wk, err := temporal_shop.NewWorker(
		ctx,
		temporal_shop.WithTemporal(clients.Temporal()),
		temporal_shop.WithInventoryClient(clients.Inventory()),
		temporal_shop.WithPubSub(clients.PubSub()),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create temporal worker: %v", err))
	}

	grpcServer, err := createGRPCServices(ctx, appCfg.GRPCConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create grpc services: %v", err))
	}
	startables := []startable{wk, grpcServer}

	for _, s := range startables {
		var current = s
		g.Go(func() error {
			if err := current.Start(ctx); err != nil {
				return err
			}
			return nil
		})
	}

	select {
	case <-quit:
		break
	case <-ctx.Done():
		break
	}

	// shutdown the things
	done()

	// limit how long we'll wait for
	timeoutCtx, timeoutCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer timeoutCancel()

	logger.Info("shutting down servers, please wait...")

	for _, s := range startables {
		s.Shutdown(timeoutCtx)
	}

	// wait for shutdown
	if err := g.Wait(); err != nil {
		panic("shutdown was not clean" + err.Error())
	}
	logger.Info("goodbye")
}
func createGRPCServices(ctx context.Context, cfg *grpc.Config) (startable, error) {
	inv, err := inventory.NewInventoryService()
	if err != nil {
		return nil, err
	}

	srv, err := grpc.NewDefaultGRPCServer(
		ctx,
		grpc.WithConfig(cfg),
		grpc.WithServices(inv),
	)
	return srv, err
}
