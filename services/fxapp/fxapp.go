package fxapp

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	StartedAt    time.Duration
	ShuttingDown = false // This variable allows to ignore some errors inherited from the shutdown
)

// Start starts application with given timeout
func Start(app *fx.App, timeout time.Duration) {
	logger := zap.L().Named("Lifecycle")
	logger.Info("๐ข Starting app...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		if graph, errGraph := fx.VisualizeError(err); errGraph == nil {
			logger.Info("๐ Error graph", zap.String("dot", graph))
		}

		logger.Fatal("๐งจ๐ฅ Unable to start app", zap.Error(err))
	}

	logger.Info("๐ App started")
}

// Shutdown closes application with given timeout
func Shutdown(app *fx.App, timeout time.Duration) {
	StartedAt = time.Now()
	ShuttingDown = true

	logger := zap.L().Named("Lifecycle")
	logger.Info("๐ข Stopping app...", zap.Duration("uptime", time.Since(appStartedAt)))

	// Listen signal to force exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Info("๐ Termination signal received", zap.String("signal", s.String()))
		os.Exit(1)
	}()

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		logger.Fatal("๐งจ๐ฅ Unable to cleanly stop app", zap.Error(err))
	}

	logger.In
}
