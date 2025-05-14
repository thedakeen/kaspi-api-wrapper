package refund

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/handlers"
	grpchandler "kaspi-api-wrapper/internal/handlers/grpc"
	refundv1 "kaspi-api-wrapper/pkg/protos/gen/go/refund"
	"log/slog"
)

type serverAPI struct {
	refundv1.UnimplementedRefundServiceServer
	log            *slog.Logger
	refundProvider handlers.RefundProvider
}

func Register(gRPC *grpc.Server, log *slog.Logger, refundProvider handlers.RefundProvider) {
	refundv1.RegisterRefundServiceServer(gRPC, &serverAPI{
		log:            log,
		refundProvider: refundProvider,
	})
}

// CreateRefundQR implements kaspiv1.RefundServiceServer
func (s *serverAPI) CreateRefundQR(ctx context.Context, req *refundv1.CreateRefundQRRequest) (*refundv1.CreateRefundQRResponse, error) {
	domainReq := domain.QRRefundCreateRequest{
		DeviceToken: req.DeviceToken,
		ExternalID:  req.ExternalId,
	}

	result, err := s.refundProvider.CreateRefundQR(ctx, domainReq)
	if err != nil {
		s.log.Error("CreateRefundQR failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundv1.CreateRefundQRResponse{
		QrToken:    result.QrToken,
		ExpireDate: timestamppb.New(result.ExpireDate),
		QrReturnId: result.QrReturnID,
		QrRefundBehaviorOptions: &refundv1.QRRefundBehaviorOptions{
			QrCodeScanEventPollingInterval: int32(result.QrRefundBehaviorOptions.QrCodeScanEventPollingInterval),
			QrCodeScanWaitTimeout:          int32(result.QrRefundBehaviorOptions.QrCodeScanWaitTimeout),
		},
	}

	return resp, nil
}

// GetRefundStatus implements kaspiv1.RefundServiceServer
func (s *serverAPI) GetRefundStatus(ctx context.Context, req *refundv1.GetRefundStatusRequest) (*refundv1.GetRefundStatusResponse, error) {
	result, err := s.refundProvider.GetRefundStatus(ctx, req.QrReturnId)
	if err != nil {
		s.log.Error("GetRefundStatus failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundv1.GetRefundStatusResponse{
		Status: result.Status,
	}

	return resp, nil
}

// GetCustomerOperations implements kaspiv1.RefundServiceServer
func (s *serverAPI) GetCustomerOperations(ctx context.Context, req *refundv1.GetCustomerOperationsRequest) (*refundv1.GetCustomerOperationsResponse, error) {
	domainReq := domain.CustomerOperationsRequest{
		DeviceToken: req.DeviceToken,
		QrReturnID:  req.QrReturnId,
		MaxResult:   req.MaxResult,
	}

	operations, err := s.refundProvider.GetCustomerOperations(ctx, domainReq)
	if err != nil {
		s.log.Error("GetCustomerOperations failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	protoOperations := make([]*refundv1.CustomerOperation, 0, len(operations))
	for _, op := range operations {
		protoOperations = append(protoOperations, &refundv1.CustomerOperation{
			QrPaymentId:     op.QrPaymentID,
			TransactionDate: timestamppb.New(op.TransactionDate),
			Amount:          op.Amount,
		})
	}

	resp := &refundv1.GetCustomerOperationsResponse{
		Operations: protoOperations,
	}

	return resp, nil
}

// GetPaymentDetails implements kaspiv1.RefundServiceServer
func (s *serverAPI) GetPaymentDetails(ctx context.Context, req *refundv1.GetPaymentDetailsRequest) (*refundv1.GetPaymentDetailsResponse, error) {
	details, err := s.refundProvider.GetPaymentDetails(ctx, req.QrPaymentId, req.DeviceToken)
	if err != nil {
		s.log.Error("GetPaymentDetails failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundv1.GetPaymentDetailsResponse{
		QrPaymentId:           details.QrPaymentID,
		TotalAmount:           details.TotalAmount,
		AvailableReturnAmount: details.AvailableReturnAmount,
		TransactionDate:       timestamppb.New(details.TransactionDate),
	}

	return resp, nil
}

// RefundPayment implements kaspiv1.RefundServiceServer
func (s *serverAPI) RefundPayment(ctx context.Context, req *refundv1.RefundPaymentRequest) (*refundv1.RefundPaymentResponse, error) {
	domainReq := domain.RefundRequest{
		DeviceToken: req.DeviceToken,
		QrPaymentID: req.QrPaymentId,
		QrReturnID:  req.QrReturnId,
		Amount:      req.Amount,
	}

	result, err := s.refundProvider.RefundPayment(ctx, domainReq)
	if err != nil {
		s.log.Error("RefundPayment failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundv1.RefundPaymentResponse{
		ReturnOperationId: result.ReturnOperationID,
	}

	return resp, nil
}
