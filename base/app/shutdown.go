package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Stopper interface {
	Close()
}

const (
	shutdownTimeout = 2 * time.Minute

	shutdownComplete       = 1
	shutdownTimeoutExpired = 2
	forcedExit             = 3
)

func (a *App) shutdown(serviceStoppers []Stopper) {
	a.logger.Info("shutting down app")

	// NOTE: make sure services are passed in correct shutdown order
	for _, s := range serviceStoppers {
		if s != nil {
			s.Close()
		}
	}

	// let the ctx chain know that the app finished
	a.Cancel()

	a.logger.Info("shutdown complete")

	// flush the log buffer
	_ = a.logger.Sync()

	os.Exit(shutdownComplete)
}

func (a *App) HandleGracefulShutdown(serviceStoppers []Stopper) {
	// create interrupt signal listener, use the app's context
	ctx, _ := signal.NotifyContext(a.ctx, syscall.SIGINT, syscall.SIGTERM)

	// block and listen for the interrupt signal
	<-ctx.Done()

	// notify user of shutdown
	a.logger.Info("shutting down gracefully, press Ctrl+C again to force")

	go func() {
		// restart listening for the force exist signal, any new context will do here
		ctx2, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-time.After(shutdownTimeout):
			a.logger.Info("shutdown timeout expired")

			os.Exit(shutdownTimeoutExpired)
		case <-ctx2.Done():
			a.logger.Info("forced exit")

			os.Exit(forcedExit)
		}
	}()

	a.shutdown(serviceStoppers)
}
