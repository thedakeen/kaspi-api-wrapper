package handlers

import (
	"context"
	"kaspi-api-wrapper/internal/domain"
)

type DeviceProvider interface {
	GetTradePoints(ctx context.Context) ([]domain.TradePoint, error)
	RegisterDevice(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDevice(ctx context.Context, deviceToken string) error
}

type DeviceEnhancedProvider interface {
	GetTradePointsEnhanced(ctx context.Context, organizationBin string) ([]domain.TradePoint, error)
	RegisterDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error
}

type PaymentProvider interface {
	CreateQR(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error)
	CreatePaymentLink(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error)
	GetPaymentStatus(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error)
}

type PaymentEnhancedProvider interface {
	CreateQREnhanced(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error)
	CreatePaymentLinkEnhanced(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error)
}

type RefundProvider interface {
	CreateRefundQR(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error)
	GetRefundStatus(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error)
	GetCustomerOperations(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error)
	GetPaymentDetails(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error)
	RefundPayment(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error)
}

type RefundEnhancedProvider interface {
	RefundPaymentEnhanced(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error)
	GetClientInfo(ctx context.Context, phoneNumber string, deviceToken int64) (*domain.ClientInfoResponse, error)
	CreateRemotePayment(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error)
	CancelRemotePayment(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error)
}

type UtilityProvider interface {
	HealthCheck(ctx context.Context) error
	TestScanQR(ctx context.Context, req domain.TestScanRequest) error
	TestConfirmPayment(ctx context.Context, req domain.TestConfirmRequest) error
	TestScanError(ctx context.Context, req domain.TestScanErrorRequest) error
	TestConfirmError(ctx context.Context, req domain.TestConfirmErrorRequest) error
}
