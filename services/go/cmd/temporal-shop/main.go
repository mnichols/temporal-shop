package main

import (
	"context"
	"github.com/temporalio/temporal-shop/services/go/pkg/clients/http"
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
	TemporalWorker *temporal.Config

	HTTPConfig *http.Config
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

	// clients
	clients := clients.MustGetClients(ctx,
		clients.WithTemporal(temporal2.NewClients(ctx,
			temporal2.WithConfig(appCfg.TemporalClient),
			temporal2.WithLogger(logger))),
	)

	defer func() {
		if perr := clients.Close(); perr != nil {
			logger.Error("failed to close clients gracefully", logur.Fields{"err": perr})
		}
	}()
	// apps

	wk, err := temporal.NewWorker(
		ctx,
		temporal.WithTemporal(clients.Temporal()),
	)
	if err != nil {
		panic("failed to create temporal worker")
	}
	startables := []startable{wk}

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
