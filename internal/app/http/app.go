package httpapp

import (
	"context"
	"errors"
	"fmt"
	httphandler "kaspi-api-wrapper/internal/handlers/http"
	"log/slog"
	"net"
	"net/http"
)

type App struct {
	log      *slog.Logger
	httpPort int
	server   *http.Server
	handlers *httphandler.Handlers
	scheme   string
}

func New(log *slog.Logger, httpPort int, handlers *httphandler.Handlers, scheme string) *App {
	return &App{
		log:      log,
		httpPort: httpPort,
		handlers: handlers,
		scheme:   scheme,
	}
}

func (app *App) Run(ctx context.Context) error {
	const op = "http.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.httpPort),
	)

	router := httphandler.NewRouter(app.log, app.handlers, app.scheme)
	r := router.Setup()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.httpPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.server = &http.Server{
		Addr:    l.Addr().String(),
		Handler: r,
	}

	log.Info("HTTP server is starting", slog.String("addr", l.Addr().String()))

	go func() {
		if err := app.server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server failed", slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (app *App) GracefulStop(ctx context.Context) error {
	const op = "http.GracefulStop"

	err := app.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to shutdown server: %w", op, err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) {
	const op = "http.Stop"

	log := app.log.With(slog.String("op", op))
	log.Info("stopping HTTP server", slog.Int("port", app.httpPort))

	err := app.GracefulStop(ctx)
	if err != nil {
		log.Error("failed to stop HTTP server", slog.String("error", err.Error()))
	}

}
