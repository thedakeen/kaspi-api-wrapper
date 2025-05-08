package main

import (
	"context"
	"kaspi-api-wrapper/internal/api/handlers"
	"kaspi-api-wrapper/internal/app"
	"kaspi-api-wrapper/internal/config"
	"kaspi-api-wrapper/internal/service"
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
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)

	log.Debug("debug enabled")

	kaspiService := service.NewKaspiService(
		log,
		cfg.KaspiAPI.Scheme,
		cfg.KaspiAPI.BaseURLBasic,
		cfg.KaspiAPI.BaseURLStd,
		cfg.KaspiAPI.BaseURLEnh,
		cfg.KaspiAPI.ApiKey,
	)

	h := handlers.NewHandlers(log, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService)

	application := app.New(log, cfg.HTTPPort, h, cfg.KaspiAPI.Scheme)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

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
