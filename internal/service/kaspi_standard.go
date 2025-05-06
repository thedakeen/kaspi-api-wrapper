package service

import (
	"context"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	"log/slog"
	"net/http"
)

//////// 	Refund service methods (standard scheme)	////////

// CreateRefundQR creates a QR code for refund (3.4.1)
func (s *KaspiService) CreateRefundQR(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error) {
	if s.scheme == "basic" {
		return nil, fmt.Errorf("refund functionality is not available in basic scheme")
	}

	const op = "service.kaspi.CreateRefundQR"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
	)

	log.Debug("creating QR token for refund")

	path := "/return/create"

	var result domain.QRRefundCreateResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("QR token for refund created successfully")

	return &result, nil
}

// GetRefundStatus gets the current status of a refund (3.4.2)
func (s *KaspiService) GetRefundStatus(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error) {
	if s.scheme == "basic" {
		return nil, fmt.Errorf("refund functionality is not available in basic scheme")
	}

	const op = "service.kaspi.GetRefundStatus"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("qrReturnID", qrReturnID),
	)

	log.Debug("getting refund status")

	path := fmt.Sprintf("/return/status/%d", qrReturnID)

	var result domain.RefundStatusResponse
	err := s.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("customer operations retrieved successfully", "status", result.Status)

	return &result, nil
}

// GetCustomerOperations gets the list of customer operations (3.4.3)
func (s *KaspiService) GetCustomerOperations(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error) {
	if s.scheme == "basic" {
		return nil, fmt.Errorf("refund functionality is not available in basic scheme")
	}

	const op = "service.kaspi.GetCustomerOperations"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Int64("qrReturnID", req.QrReturnID),
	)

	log.Debug("getting customer operations")

	path := "/return/operations"

	var result []domain.CustomerOperation
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("customer operations retrieved successfully", "count", len(result))

	return result, nil
}

// GetPaymentDetails gets the details of a payment (3.4.4)
func (s *KaspiService) GetPaymentDetails(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error) {
	if s.scheme == "basic" {
		return nil, fmt.Errorf("refund functionality is not available in basic scheme")
	}

	const op = "service.kaspi.GetPaymentDetails"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("qrPaymentID", qrPaymentID),
		slog.String("deviceToken", deviceToken),
	)

	log.Debug("getting payment details")

	path := fmt.Sprintf("/payment/details?QrPaymentId=%d&DeviceToken=%s", qrPaymentID, deviceToken)

	var result domain.PaymentDetailsResponse
	err := s.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment details retrieved successfully")

	return &result, nil
}

// RefundPayment initiates a payment refund (3.4.5)
func (s *KaspiService) RefundPayment(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error) {
	if s.scheme == "basic" {
		return nil, fmt.Errorf("refund functionality is not available in basic scheme")
	}

	const op = "service.kaspi.RefundPayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Int64("qrPaymentID", req.QrPaymentID),
		slog.Float64("amount", req.Amount),
	)

	log.Debug("initiating payment refund")

	path := "/payment/return"

	var result domain.RefundResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment refund initiated successfully")

	return &result, nil
}

//////// 	End of refund service methods (standard scheme)	////////
