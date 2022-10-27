package main

import (
	"context"
	"github.com/temporalio/temporal-shop/web/bff/build"
	"github.com/temporalio/temporal-shop/web/bff/config"
	"github.com/temporalio/temporal-shop/web/bff/internal/clients"
	temporalClient "github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	httpServer "github.com/temporalio/temporal-shop/web/bff/internal/http/server"
	"github.com/temporalio/temporal-shop/web/bff/internal/instrumentation/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"logur.dev/logur"
)

type startable interface {
	Start(context.Context) error
	Shutdown(context.Context)
}

type appConfig struct {
	Log            *log.Config
	HttpServer     *httpServer.Config
	TemporalClient *temporalClient.Config
}

func main() {
	ctx := context.Background()
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
	logger = log.WithFields(logger, log.Fields{
		"version":    build.Version,
		"build_date": build.BuildDate,
		"commit":     build.Commit,
	})
	ctx = log.WithLogger(ctx, logger)

	logger.Info("config", log.Fields{"cfg": appCfg})
	// clients
	clients := clients.MustGetClients(ctx,
		clients.WithTemporal(temporalClient.NewClients(ctx,
			temporalClient.WithConfig(appCfg.TemporalClient),
			temporalClient.WithLogger(logger))),
	)
	defer func() {
		if perr := clients.Close(); perr != nil {
			logger.Error("failed to close clients gracefully", logur.Fields{"err": perr})
		}
	}()
	// apps
	hs, err := httpServer.NewServer(ctx,
		httpServer.WithConfig(appCfg.HttpServer),
		httpServer.WithTemporalClients(clients.Temporal()),
		httpServer.WithLogger(logger),
	)
	if err != nil {
		panic("failed to create http server: " + err.Error())
	}

	startables := []startable{hs}

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
