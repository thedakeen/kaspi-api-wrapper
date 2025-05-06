package service

import (
	"context"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	"log/slog"
	"net/http"
	"net/url"
)

//////// 	Refund service	methods	(enhanced) 	////////

// RefundPaymentEnhanced initiates a payment refund without customer participation (4.5)
func (s *KaspiService) RefundPaymentEnhanced(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error) {
	if s.scheme != "enhanced" {
		return nil, fmt.Errorf("enhanced refund functionality is only available in enhanced scheme")
	}

	const op = "service.kaspi.RefundPaymentEnhanced"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Int64("qrPaymentID", req.QrPaymentID),
		slog.Float64("amount", req.Amount),
	)

	log.Debug("initiating payment refund (enhanced)")

	path := "/payment/return"

	var result domain.RefundResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment refund initiated successfully")

	return &result, nil
}

// GetClientInfo retrieves client information by phone number (4.6.1)
func (s *KaspiService) GetClientInfo(ctx context.Context, phoneNumber, deviceToken string) (*domain.ClientInfoResponse, error) {
	if s.scheme != "enhanced" {
		return nil, fmt.Errorf("remote payment functionality is only available in enhanced scheme")
	}

	const op = "service.kaspi.GetClientInfo"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", deviceToken),
	)

	log.Debug("getting client information by phone number")

	path := fmt.Sprintf("/remote/client-info?phoneNumber=%s&deviceToken=%s",
		url.QueryEscape(phoneNumber), url.QueryEscape(deviceToken))

	var result domain.ClientInfoResponse
	err := s.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("client information retrieved successfully")

	return &result, nil
}

// CreateRemotePayment creates a remote payment request (4.6.2)
func (s *KaspiService) CreateRemotePayment(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error) {
	if s.scheme != "enhanced" {
		return nil, fmt.Errorf("remote payment functionality is only available in enhanced scheme")
	}

	const op = "service.kaspi.CreateRemotePayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
	)

	log.Debug("creating remote payment request")

	path := "/remote/create"

	var result domain.RemotePaymentResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("remote payment request created successfully")

	return &result, nil
}

// CancelRemotePayment cancels a remote payment request (4.6.3)
func (s *KaspiService) CancelRemotePayment(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error) {
	if s.scheme != "enhanced" {
		return nil, fmt.Errorf("remote payment functionality is only available in enhanced scheme")
	}

	const op = "service.kaspi.CancelRemotePayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Int64("qrPaymentID", req.QrPaymentID),
	)

	log.Debug("canceling remote payment request")

	path := "/remote/cancel"

	var result domain.RemotePaymentCancelResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("remote payment request canceled successfully")

	return &result, nil
}

//////// 	End of refund service	methods	(enhanced) 	////////
