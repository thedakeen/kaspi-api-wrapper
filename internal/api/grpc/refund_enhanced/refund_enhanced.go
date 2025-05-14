package refund_enhanced

import (
	"context"
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/domain"
	refundenhancedv1 "kaspi-api-wrapper/pkg/protos/gen/go/refund_enhanced"
	"log/slog"
	"strconv"
)

type serverAPI struct {
	refundenhancedv1.UnimplementedEnhancedRefundServiceServer
	log                    *slog.Logger
	refundEnhancedProvider api.RefundEnhancedProvider
}

func Register(gRPC *grpc.Server, log *slog.Logger, refundEnhancedProvider api.RefundEnhancedProvider) {
	refundenhancedv1.RegisterEnhancedRefundServiceServer(gRPC, &serverAPI{
		log:                    log,
		refundEnhancedProvider: refundEnhancedProvider,
	})
}

// RefundPaymentEnhanced implements kaspiv1.EnhancedRefundServiceServer
func (s *serverAPI) RefundPaymentEnhanced(ctx context.Context, req *refundenhancedv1.RefundPaymentEnhancedRequest) (*refundenhancedv1.RefundPaymentEnhancedResponse, error) {
	domainReq := domain.EnhancedRefundRequest{
		DeviceToken:     req.DeviceToken,
		QrPaymentID:     req.QrPaymentId,
		Amount:          req.Amount,
		OrganizationBin: req.OrganizationBin,
	}

	result, err := s.refundEnhancedProvider.RefundPaymentEnhanced(ctx, domainReq)
	if err != nil {
		s.log.Error("RefundPaymentEnhanced failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundenhancedv1.RefundPaymentEnhancedResponse{
		ReturnOperationId: result.ReturnOperationID,
	}

	return resp, nil
}

// GetClientInfo implements kaspiv1.EnhancedRefundServiceServer
func (s *serverAPI) GetClientInfo(ctx context.Context, req *refundenhancedv1.GetClientInfoRequest) (*refundenhancedv1.GetClientInfoResponse, error) {
	info, err := s.refundEnhancedProvider.GetClientInfo(ctx, req.PhoneNumber, req.DeviceToken)
	if err != nil {
		s.log.Error("GetClientInfo failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundenhancedv1.GetClientInfoResponse{
		ClientName: info.ClientName,
	}

	return resp, nil
}

// CreateRemotePayment implements kaspiv1.EnhancedRefundServiceServer
func (s *serverAPI) CreateRemotePayment(ctx context.Context, req *refundenhancedv1.CreateRemotePaymentRequest) (*refundenhancedv1.CreateRemotePaymentResponse, error) {
	deviceToken, err := strconv.ParseInt(req.DeviceToken, 10, 64)
	if err != nil {
		return nil, grpchandler.HandleError(err, s.log)
	}

	domainReq := domain.RemotePaymentRequest{
		OrganizationBin: req.OrganizationBin,
		Amount:          req.Amount,
		PhoneNumber:     req.PhoneNumber,
		DeviceToken:     deviceToken,
		Comment:         req.Comment,
	}

	result, err := s.refundEnhancedProvider.CreateRemotePayment(ctx, domainReq)
	if err != nil {
		s.log.Error("CreateRemotePayment failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundenhancedv1.CreateRemotePaymentResponse{
		QrPaymentId: result.QrPaymentID,
	}

	return resp, nil
}

// CancelRemotePayment implements kaspiv1.EnhancedRefundServiceServer
func (s *serverAPI) CancelRemotePayment(ctx context.Context, req *refundenhancedv1.CancelRemotePaymentRequest) (*refundenhancedv1.CancelRemotePaymentResponse, error) {
	domainReq := domain.RemotePaymentCancelRequest{
		OrganizationBin: req.OrganizationBin,
		QrPaymentID:     req.QrPaymentId,
		DeviceToken:     req.DeviceToken,
	}

	result, err := s.refundEnhancedProvider.CancelRemotePayment(ctx, domainReq)
	if err != nil {
		s.log.Error("CancelRemotePayment failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	resp := &refundenhancedv1.CancelRemotePaymentResponse{
		Status: result.Status,
	}

	return resp, nil
}
