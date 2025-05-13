package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/api/grpc/device"
	"kaspi-api-wrapper/internal/api/grpc/payment"
	"kaspi-api-wrapper/internal/api/grpc/refund"
	"kaspi-api-wrapper/internal/api/grpc/refund_enhanced"
	"kaspi-api-wrapper/internal/api/grpc/utility"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	grpcPort   int
}

func New(log *slog.Logger, grpcPort int, handlers *grpchandler.Handlers,
) *App {
	gRPCServer := grpc.NewServer()

	device.Register(gRPCServer, handlers.DeviceProvider, handlers.DeviceEnhancedProvider)
	payment.Register(gRPCServer, handlers.PaymentProvider, handlers.PaymentEnhancedProvider)
	refund.Register(gRPCServer, handlers.RefundProvider)
	refund_enhanced.Register(gRPCServer, handlers.RefundEnhancedProvider)
	utility.Register(gRPCServer, handlers.UtilityProvider)

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
