package app

import (
	"kaspi-api-wrapper/internal/api/http"
	grpcapp "kaspi-api-wrapper/internal/app/grpc"
	"kaspi-api-wrapper/internal/app/http"
	"kaspi-api-wrapper/internal/service"
	"log/slog"
)

type App struct {
	HTTPSrv  *httpapp.App
	GRPCSrv  *grpcapp.App
	Handlers *http.Handlers
}

func New(log *slog.Logger, httpPort int, handlers *http.Handlers, scheme string, grpcPort int, kaspiService *service.KaspiService) *App {
	httpApp := httpapp.New(log, httpPort, handlers, scheme)
	grcpApp := grpcapp.New(log, grpcPort, kaspiService)

	return &App{
		httpApp,
		grcpApp,
		handlers,
	}
}
