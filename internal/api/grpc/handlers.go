package grpchandler

import (
	"kaspi-api-wrapper/internal/api"
	"log/slog"
)

// Handlers contains all gRPC handlers for the API
type Handlers struct {
	log             *slog.Logger
	deviceProvider  api.DeviceProvider
	paymentProvider api.PaymentProvider
	utilityProvider api.UtilityProvider
	refundProvider  api.RefundProvider

	deviceEnhancedProvider  api.DeviceEnhancedProvider
	paymentEnhancedProvider api.PaymentEnhancedProvider
	refundEnhancedProvider  api.RefundEnhancedProvider
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
		deviceProvider:  deviceProvider,
		paymentProvider: paymentProvider,
		utilityProvider: utilityProvider,
		refundProvider:  refundProvider,

		deviceEnhancedProvider:  deviceEnhancedProvider,
		paymentEnhancedProvider: paymentEnhancedProvider,
		refundEnhancedProvider:  refundEnhancedProvider,
		//kaspiSvc: kaspiSvc,
	}
}
