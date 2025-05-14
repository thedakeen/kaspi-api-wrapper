package utility

import (
	"context"
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/handlers"
	grpchandler "kaspi-api-wrapper/internal/handlers/grpc"
	utilityv1 "kaspi-api-wrapper/pkg/protos/gen/go/utility"
	"log/slog"
)

type serverAPI struct {
	utilityv1.UnimplementedUtilityServiceServer
	log             *slog.Logger
	utilityProvider handlers.UtilityProvider
}

func Register(gRPC *grpc.Server, log *slog.Logger, utilityProvider handlers.UtilityProvider) {
	utilityv1.RegisterUtilityServiceServer(gRPC, &serverAPI{
		log:             log,
		utilityProvider: utilityProvider,
	})
}

func RegisterTest(log *slog.Logger, utilityProvider handlers.UtilityProvider) utilityv1.UtilityServiceServer {
	return &serverAPI{
		log:             log,
		utilityProvider: utilityProvider,
	}
}

// HealthCheck implements kaspiv1.UtilityServiceServer
func (s *serverAPI) HealthCheck(ctx context.Context, req *utilityv1.HealthCheckRequest) (*utilityv1.HealthCheckResponse, error) {
	err := s.utilityProvider.HealthCheck(ctx)
	if err != nil {
		s.log.Error("HealthCheck failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &utilityv1.HealthCheckResponse{
		Status: "OK",
	}, nil
}

// TestScanQR implements kaspiv1.UtilityServiceServer
func (s *serverAPI) TestScanQR(ctx context.Context, req *utilityv1.TestScanQRRequest) (*utilityv1.TestScanQRResponse, error) {
	domainReq := domain.TestScanRequest{
		QrPaymentID: req.QrPaymentId,
	}

	err := s.utilityProvider.TestScanQR(ctx, domainReq)
	if err != nil {
		s.log.Error("TestScanQR failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &utilityv1.TestScanQRResponse{
		Message: "QR scan simulation successful",
	}, nil
}

// TestConfirmPayment implements kaspiv1.UtilityServiceServer
func (s *serverAPI) TestConfirmPayment(ctx context.Context, req *utilityv1.TestConfirmPaymentRequest) (*utilityv1.TestConfirmPaymentResponse, error) {
	domainReq := domain.TestConfirmRequest{
		QrPaymentID: req.QrPaymentId,
	}

	err := s.utilityProvider.TestConfirmPayment(ctx, domainReq)
	if err != nil {
		s.log.Error("TestConfirmPayment failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &utilityv1.TestConfirmPaymentResponse{
		Message: "Payment confirmation simulation successful",
	}, nil
}

// TestScanError implements kaspiv1.UtilityServiceServer
func (s *serverAPI) TestScanError(ctx context.Context, req *utilityv1.TestScanErrorRequest) (*utilityv1.TestScanErrorResponse, error) {
	domainReq := domain.TestScanErrorRequest{
		QrPaymentID: req.QrPaymentId,
	}

	err := s.utilityProvider.TestScanError(ctx, domainReq)
	if err != nil {
		s.log.Error("TestScanError failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &utilityv1.TestScanErrorResponse{
		Message: "QR scan error simulation successful",
	}, nil
}

// TestConfirmError implements kaspiv1.UtilityServiceServer
func (s *serverAPI) TestConfirmError(ctx context.Context, req *utilityv1.TestConfirmErrorRequest) (*utilityv1.TestConfirmErrorResponse, error) {
	domainReq := domain.TestConfirmErrorRequest{
		QrPaymentID: req.QrPaymentId,
	}

	err := s.utilityProvider.TestConfirmError(ctx, domainReq)
	if err != nil {
		s.log.Error("TestConfirmError failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &utilityv1.TestConfirmErrorResponse{
		Message: "Payment confirmation error simulation successful",
	}, nil
}
