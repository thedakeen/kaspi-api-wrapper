package app

import (
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/api/http"
	grpcapp "kaspi-api-wrapper/internal/app/grpc"
	"kaspi-api-wrapper/internal/app/http"
	"kaspi-api-wrapper/internal/service"
	"log/slog"
)

type App struct {
	HTTPSrv      *httpapp.App
	GRPCSrv      *grpcapp.App
	httpHandlers *http.Handlers
	grpcHandlers *grpchandler.Handlers
}

func New(log *slog.Logger, httpPort int, scheme string, grpcPort int, kaspiService *service.KaspiService) *App {
	httpHandlers := http.NewHandlers(log, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService)
	grpcHandlers := grpchandler.NewHandlers(log, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService, kaspiService)

	httpApp := httpapp.New(log, httpPort, httpHandlers, scheme)
	grpcApp := grpcapp.New(log, grpcPort, grpcHandlers, scheme)

	return &App{
		httpApp,
		grpcApp,
		httpHandlers,
		grpcHandlers,
	}
}
