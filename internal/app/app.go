package app

import (
	"kaspi-api-wrapper/internal/api/handlers"
	"kaspi-api-wrapper/internal/app/http"
	"log/slog"
)

type App struct {
	HTTPSrv  *http.App
	Handlers *handlers.Handlers
}

func New(log *slog.Logger, httpPort int, handlers *handlers.Handlers) *App {
	httpApp := http.New(log, httpPort, handlers)

	return &App{
		httpApp,
		handlers,
	}
}
