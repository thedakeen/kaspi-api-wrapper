package payment

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kaspi-api-wrapper/internal/api"
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/validator"
	paymentv1 "kaspi-api-wrapper/pkg/protos/gen/go/payment"
	"log/slog"
)

type serverAPI struct {
	paymentv1.UnimplementedPaymentServiceServer
	log                     *slog.Logger
	paymentProvider         api.PaymentProvider
	paymentEnhancedProvider api.PaymentEnhancedProvider
}

func Register(gRPC *grpc.Server, log *slog.Logger, paymentProvider api.PaymentProvider, paymentEnhancedProvider api.PaymentEnhancedProvider) {
	paymentv1.RegisterPaymentServiceServer(gRPC, &serverAPI{
		log:                     log,
		paymentProvider:         paymentProvider,
		paymentEnhancedProvider: paymentEnhancedProvider,
	})
}

// CreateQR implements kaspiv1.PaymentServiceServer
func (s *serverAPI) CreateQR(ctx context.Context, req *paymentv1.CreateQRRequest) (*paymentv1.CreateQRResponse, error) {
	domainReq := domain.QRCreateRequest{
		DeviceToken: req.DeviceToken,
		Amount:      req.Amount,
		ExternalID:  req.ExternalId,
	}

	if err := validator.ValidateQRCreateRequest(domainReq); err != nil {
		return nil, validator.GRPCError(err)
	}

	result, err := s.paymentProvider.CreateQR(ctx, domainReq)
	if err != nil {
		s.log.Error("CreateQR failed", "error", err.Error())
		return nil, grpchandler.HandleKaspiError(err, s.log)
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

// CreatePaymentLink implements kaspiv1.PaymentServiceServer
func (s *serverAPI) CreatePaymentLink(ctx context.Context, req *paymentv1.CreatePaymentLinkRequest) (*paymentv1.CreatePaymentLinkResponse, error) {
	domainReq := domain.PaymentLinkCreateRequest{
		DeviceToken: req.DeviceToken,
		Amount:      req.Amount,
		ExternalID:  req.ExternalId,
	}

	if err := validator.ValidatePaymentLinkCreateRequest(domainReq); err != nil {
		return nil, validator.GRPCError(err)
	}

	result, err := s.paymentProvider.CreatePaymentLink(ctx, domainReq)
	if err != nil {
		s.log.Error("CreatePaymentLink failed", "error", err.Error())
		return nil, grpchandler.HandleKaspiError(err, s.log)
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

// GetPaymentStatus implements kaspiv1.PaymentServiceServer
func (s *serverAPI) GetPaymentStatus(ctx context.Context, req *paymentv1.GetPaymentStatusRequest) (*paymentv1.GetPaymentStatusResponse, error) {
	if req.QrPaymentId <= 0 {
		return nil, validator.GRPCError(&validator.ValidationError{
			Field:   "qrPaymentId",
			Message: "payment ID must be a positive number",
			Err:     validator.ErrInvalidID,
		})
	}

	result, err := s.paymentProvider.GetPaymentStatus(ctx, req.QrPaymentId)
	if err != nil {
		s.log.Error("GetPaymentStatus failed", "error", err.Error())
		return nil, grpchandler.HandleKaspiError(err, s.log)
	}

	resp := &paymentv1.GetPaymentStatusResponse{
		Status:        result.Status,
		TransactionId: result.TransactionID,
		LoanOfferName: result.LoanOfferName,
		// TODO: by docs LoanTerm is just int, maybe mistake
		LoanTerm:    int64(result.LoanTerm),
		IsOffer:     result.IsOffer,
		ProductType: result.ProductType,
		Amount:      result.Amount,
		StoreName:   result.StoreName,
		Address:     result.Address,
		City:        result.City,
	}

	return resp, nil
}
