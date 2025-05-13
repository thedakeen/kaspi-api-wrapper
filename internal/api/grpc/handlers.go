package grpchandler

import (
	"kaspi-api-wrapper/internal/api"
	"log/slog"
)

// Handlers contains all gRPC handlers for the API
type Handlers struct {
	log             *slog.Logger
	DeviceProvider  api.DeviceProvider
	PaymentProvider api.PaymentProvider
	UtilityProvider api.UtilityProvider
	RefundProvider  api.RefundProvider

	DeviceEnhancedProvider  api.DeviceEnhancedProvider
	PaymentEnhancedProvider api.PaymentEnhancedProvider
	RefundEnhancedProvider  api.RefundEnhancedProvider
	//kaspiSvc *service.KaspiService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(
	log *slog.Logger,
	deviceProvider api.DeviceProvider,
	paymentProvider api.PaymentProvider,
	utilityProvider api.UtilityProvider,
	refundProvider api.RefundProvider,

	deviceEnhancedProvider api.DeviceEnhancedProvider,
	paymentEnhancedProvider api.PaymentEnhancedProvider,
	refundEnhancedProvider api.RefundEnhancedProvider,
) *Handlers {
	return &Handlers{
		log:             log,
		DeviceProvider:  deviceProvider,
		PaymentProvider: paymentProvider,
		UtilityProvider: utilityProvider,
		RefundProvider:  refundProvider,

		DeviceEnhancedProvider:  deviceEnhancedProvider,
		PaymentEnhancedProvider: paymentEnhancedProvider,
		RefundEnhancedProvider:  refundEnhancedProvider,
		//kaspiSvc: kaspiSvc,
	}
}
