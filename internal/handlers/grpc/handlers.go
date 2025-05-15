package grpchandler

import (
	"kaspi-api-wrapper/internal/handlers"
	"log/slog"
)

// Handlers contains all gRPC handlers for the API
type Handlers struct {
	log             *slog.Logger
	DeviceProvider  handlers.DeviceProvider
	PaymentProvider handlers.PaymentProvider
	UtilityProvider handlers.UtilityProvider
	RefundProvider  handlers.RefundProvider

	DeviceEnhancedProvider  handlers.DeviceEnhancedProvider
	PaymentEnhancedProvider handlers.PaymentEnhancedProvider
	RefundEnhancedProvider  handlers.RefundEnhancedProvider
	//kaspiSvc *service.KaspiService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(
	log *slog.Logger,
	deviceProvider handlers.DeviceProvider,
	paymentProvider handlers.PaymentProvider,
	utilityProvider handlers.UtilityProvider,
	refundProvider handlers.RefundProvider,

	deviceEnhancedProvider handlers.DeviceEnhancedProvider,
	paymentEnhancedProvider handlers.PaymentEnhancedProvider,
	refundEnhancedProvider handlers.RefundEnhancedProvider,
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
