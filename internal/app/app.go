package app

import (
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/api/http"
	grpcapp "kaspi-api-wrapper/internal/app/grpc"
	"kaspi-api-wrapper/internal/app/http"
	"log/slog"
)

type App struct {
	HTTPSrv      *httpapp.App
	GRPCSrv      *grpcapp.App
	httpHandlers *http.Handlers
	grpcHandlers *grpchandler.Handlers
}

func New(log *slog.Logger, httpPort int, httpHandlers *http.Handlers, scheme string, grpcPort int, grpcHandlers *grpchandler.Handlers) *App {
	httpApp := httpapp.New(log, httpPort, httpHandlers, scheme)
	grpcApp := grpcapp.New(log, grpcPort, grpcHandlers)

	return &App{
		httpApp,
		grpcApp,
		httpHandlers,
		grpcHandlers,
	}
}
