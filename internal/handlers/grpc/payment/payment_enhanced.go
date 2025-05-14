package payment

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kaspi-api-wrapper/internal/domain"
	grpchandler "kaspi-api-wrapper/internal/handlers/grpc"
	paymentv1 "kaspi-api-wrapper/pkg/protos/gen/go/payment"
)

// CreateQREnhanced implements kaspiv1.PaymentServiceServer
func (s *serverAPI) CreateQREnhanced(ctx context.Context, req *paymentv1.CreateQREnhancedRequest) (*paymentv1.CreateQRResponse, error) {
	domainReq := domain.EnhancedQRCreateRequest{
		DeviceToken:     req.DeviceToken,
		Amount:          req.Amount,
		ExternalID:      req.ExternalId,
		OrganizationBin: req.OrganizationBin,
	}

	result, err := s.paymentEnhancedProvider.CreateQREnhanced(ctx, domainReq)
	if err != nil {
		// Log only errors
		s.log.Error("CreateQREnhanced failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &paymentv1.CreateQRResponse{
		QrToken:        result.QrToken,
		ExpireDate:     timestamppb.New(result.ExpireDate),
		QrPaymentId:    result.QrPaymentID,
		PaymentMethods: result.PaymentMethods,
		QrPaymentBehaviorOptions: &paymentv1.QRPaymentBehaviorOptions{
			StatusPollingInterval:      int64(result.QrPaymentBehaviorOptions.StatusPollingInterval),
			QrCodeScanWaitTimeout:      int64(result.QrPaymentBehaviorOptions.QrCodeScanWaitTimeout),
			PaymentConfirmationTimeout: int64(result.QrPaymentBehaviorOptions.PaymentConfirmationTimeout),
		},
	}

	return resp, nil
}

// CreatePaymentLinkEnhanced implements kaspiv1.PaymentServiceServer
func (s *serverAPI) CreatePaymentLinkEnhanced(ctx context.Context, req *paymentv1.CreatePaymentLinkEnhancedRequest) (*paymentv1.CreatePaymentLinkResponse, error) {
	domainReq := domain.EnhancedPaymentLinkCreateRequest{
		DeviceToken:     req.DeviceToken,
		Amount:          req.Amount,
		ExternalID:      req.ExternalId,
		OrganizationBin: req.OrganizationBin,
	}

	result, err := s.paymentEnhancedProvider.CreatePaymentLinkEnhanced(ctx, domainReq)
	if err != nil {
		s.log.Error("CreatePaymentLinkEnhanced failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &paymentv1.CreatePaymentLinkResponse{
		PaymentLink:    result.PaymentLink,
		ExpireDate:     timestamppb.New(result.ExpireDate),
		PaymentId:      result.PaymentID,
		PaymentMethods: result.PaymentMethods,
		PaymentBehaviorOptions: &paymentv1.PaymentBehaviorOptions{
			StatusPollingInterval:      int64(result.PaymentBehaviorOptions.StatusPollingInterval),
			LinkActivationWaitTimeout:  int64(result.PaymentBehaviorOptions.LinkActivationWaitTimeout),
			PaymentConfirmationTimeout: int64(result.PaymentBehaviorOptions.PaymentConfirmationTimeout),
		},
	}

	return resp, nil
}
