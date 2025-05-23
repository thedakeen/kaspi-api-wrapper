package main

import (
	"context"
	"fmt"
	"kaspi-api-wrapper/internal/app"
	"kaspi-api-wrapper/internal/config"
	"kaspi-api-wrapper/internal/service"
	"kaspi-api-wrapper/internal/storage/postgres"
	"kaspi-api-wrapper/pkg/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	wg sync.WaitGroup
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(2)

	log.Debug("debug enabled")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode,
	)

	log.Info("connecting to database", "host", cfg.Database.Host, "dbname", cfg.Database.Name)
	storage, err := postgres.New(dsn)
	if err != nil {
		panic(err)
	}
	defer storage.Stop()

	tlsConfig := &service.TLSConfig{
		Password:      cfg.KaspiAPI.KeyPass,
		PfxFile:       cfg.KaspiAPI.PfxFile,
		RootCAFile:    cfg.KaspiAPI.RootCAFile,
		UseClientCert: true,
	}

	kaspiService := service.NewKaspiService(
		log,
		cfg.KaspiAPI.Scheme,
		cfg.KaspiAPI.BaseURLBasic,
		cfg.KaspiAPI.BaseURLStd,
		cfg.KaspiAPI.BaseURLEnh,
		cfg.KaspiAPI.ApiKey,

		tlsConfig,

		storage,
	)

	application := app.New(log, cfg.HTTPPort, cfg.KaspiAPI.Scheme, cfg.GRPCPort, kaspiService)

	go func() {
		defer wg.Done()
		if err := application.HTTPSrv.Run(ctx); err != nil {
			log.Error("failed to start application", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		if err := application.GRPCSrv.Run(); err != nil {
			log.Error("failed to start application", "error", err)
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	log.Info("application started")

	// Wait for interrupt signal
	<-shutdown

	log.Info("shutting down application...")

	// Stop application
	application.HTTPSrv.Stop(ctx)
	application.GRPCSrv.Stop()

	wg.Wait()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
