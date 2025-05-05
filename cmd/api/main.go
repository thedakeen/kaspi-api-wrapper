package main

import (
	"context"
	"kaspi-api-wrapper/internal/api/handlers"
	"kaspi-api-wrapper/internal/app"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	wg sync.WaitGroup
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)

	log.Debug("debug enabled")

	h := handlers.NewHandlers(log)

	application := app.New(log, 8080, h)

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start application in a goroutine
	go func() {
		defer wg.Done()
		if err := application.HTTPSrv.Run(ctx); err != nil {
			log.Error("failed to start application", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("application started")

	// Wait for interrupt signal
	<-shutdown

	log.Info("shutting down application...")

	// Stop application
	application.HTTPSrv.Stop(ctx)

	wg.Wait()

	log.Info("application stopped")
}
