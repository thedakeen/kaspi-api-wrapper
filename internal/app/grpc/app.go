package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	grpchandler "kaspi-api-wrapper/internal/handlers/grpc"
	"kaspi-api-wrapper/internal/handlers/grpc/device"
	grpcmiddleware "kaspi-api-wrapper/internal/handlers/grpc/middleware"
	"kaspi-api-wrapper/internal/handlers/grpc/payment"
	"kaspi-api-wrapper/internal/handlers/grpc/refund"
	"kaspi-api-wrapper/internal/handlers/grpc/refund_enhanced"
	"kaspi-api-wrapper/internal/handlers/grpc/utility"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	grpcPort   int
}

func New(log *slog.Logger, grpcPort int, handlers *grpchandler.Handlers, scheme string) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.SchemeInterceptor(scheme)))

	device.Register(gRPCServer, log, handlers.DeviceProvider, handlers.DeviceEnhancedProvider)
	payment.Register(gRPCServer, log, handlers.PaymentProvider, handlers.PaymentEnhancedProvider)
	refund.Register(gRPCServer, log, handlers.RefundProvider)
	refund_enhanced.Register(gRPCServer, log, handlers.RefundEnhancedProvider)
	utility.Register(gRPCServer, log, handlers.UtilityProvider)

	return &App{
		log:        log,
		grpcPort:   grpcPort,
		gRPCServer: gRPCServer,
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.grpcPort),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.grpcPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is starting", slog.String("addr", l.Addr().String()))

	err = app.gRPCServer.Serve(l)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", app.grpcPort))

	app.gRPCServer.GracefulStop()
}
