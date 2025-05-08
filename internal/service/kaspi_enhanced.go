package service

import (
	"context"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	"log/slog"
	"net/http"
	"net/url"
)

/*
Отдельные Enhanced методы для третьей схемы интеграции (вариант с усиленной безопасностью)
отличаются тем, что:

1. Требуют дополнительный параметр OrganizationBin для идентификации организации
2. Работают через отдельный базовый URL (mtokentest.kaspi.kz:8545)
3. Некоторые методы (e.g. 4.5) имеют другие входные параметры

Таким образом было решено создать отдельные методы, вдобавок это добавляет:
- Четкое разделение между разными схемами интеграции
- Предотвращение ошибок при неправильном использовании методов
*/

//////// 	Device service	methods	(enhanced) 	////////

// GetTradePointsEnhanced gets a list of trade points in the enhanced scheme (4.2.2)
func (s *KaspiService) GetTradePointsEnhanced(ctx context.Context, organizationBin string) ([]domain.TradePoint, error) {
	const op = "service.kaspi.GetTradePointsEnhanced"

	log := s.log.With(
		slog.String("op", op),
		slog.String("organizationBin", organizationBin),
	)

	log.Debug("getting trade points (enhanced)")

	path := fmt.Sprintf("/partner/tradepoints/%s", organizationBin)

	var result []domain.TradePoint
	err := s.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("trade points retrieved successfully (enhanced)")

	return result, nil
}

// RegisterDeviceEnhanced registers a device in the enhanced scheme (4.2.3)
func (s *KaspiService) RegisterDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
	const op = "service.kaspi.RegisterDeviceEnhanced"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceID", req.DeviceID),
		slog.Int64("tradePointID", req.TradePointID),
		slog.String("organizationBin", req.OrganizationBin),
	)

	log.Debug("registering device (enhanced)")

	path := "/device/register"

	var result domain.DeviceRegisterResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("device registered successfully (enhanced)")

	return &result, nil
}

// DeleteDeviceEnhanced deletes a device in the enhanced scheme (4.2.4)
func (s *KaspiService) DeleteDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error {
	const op = "service.kaspi.DeleteDeviceEnhanced"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.String("organizationBin", req.OrganizationBin),
	)

	log.Debug("deleting device (enhanced)")

	path := "/device/delete"

	err := s.request(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("device deleted successfully (enhanced)")

	return nil
}

//////// 	End of device service	methods	(enhanced) 	////////

//////// 	Payment service	methods	(enhanced) 	////////

// CreateQREnhanced creates a QR code for payment in the enhanced scheme (4.3.1)
func (s *KaspiService) CreateQREnhanced(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error) {
	const op = "service.kaspi.CreateQREnhanced"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
		slog.String("organizationBin", req.OrganizationBin),
	)

	log.Debug("creating QR (enhanced)")

	path := "/qr/create"

	var result domain.QRCreateResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("QR created successfully (enhanced)")

	return &result, nil
}

// CreatePaymentLinkEnhanced creates a payment link in the enhanced scheme (4.3.2)
func (s *KaspiService) CreatePaymentLinkEnhanced(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
	const op = "service.kaspi.CreatePaymentLinkEnhanced"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
		slog.String("organizationBin", req.OrganizationBin),
	)

	log.Debug("creating payment link (enhanced)")

	path := "/qr/create-link"

	var result domain.PaymentLinkCreateResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment link created successfully (enhanced)")

	return &result, nil
}

//////// 	End of payment service	methods	(enhanced) 	////////

//////// 	Refund service	methods	(enhanced) 	////////

// RefundPaymentEnhanced initiates a payment refund without customer participation (4.5)
func (s *KaspiService) RefundPaymentEnhanced(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error) {
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
