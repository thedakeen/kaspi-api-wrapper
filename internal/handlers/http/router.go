package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	middleware2 "kaspi-api-wrapper/internal/handlers/http/middleware"
	"log/slog"
)

type Router struct {
	log      *slog.Logger
	handlers *Handlers
	scheme   string
}

func NewRouter(log *slog.Logger, handlers *Handlers, scheme string) *Router {
	return &Router{
		log:      log,
		handlers: handlers,
		scheme:   scheme,
	}
}

func (r *Router) Setup() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware2.Logger(r.log))
	router.Use(middleware.Recoverer)

	router.Get("/health", r.handlers.HealthCheck)

	router.Route("/api", func(apiRouter chi.Router) {
		// 2.2.2 - Get trade points
		apiRouter.Get("/tradepoints", r.handlers.GetTradePoints)

		// 2.2.3 - Register device
		apiRouter.Post("/device/register", r.handlers.RegisterDevice)

		// 2.2.4 - Delete device
		apiRouter.Post("/device/delete", r.handlers.DeleteDevice)

		// 2.3.1 - Create QR code
		apiRouter.Post("/qr/create", r.handlers.CreateQR)

		// 2.3.2 - Create payment link
		apiRouter.Post("/qr/create-link", r.handlers.CreatePaymentLink)

		// 2.3.3 - Get payment status
		apiRouter.Get("/payment/status/{qrPaymentId}", r.handlers.GetPaymentStatus)

		// Standard scheme endpoints (available in standard and enhanced schemes)
		standardScheme := middleware2.SchemeMiddleware(r.scheme, "standard")

		// 3.4.1 - Create refund QR code
		apiRouter.With(standardScheme).Post("/return/create", r.handlers.CreateRefundQR)

		// 3.4.2 - Get refund status
		apiRouter.With(standardScheme).Get("/return/status/{qrReturnId}", r.handlers.GetRefundStatus)

		// 3.4.3 - Get customer operations
		apiRouter.With(standardScheme).Post("/return/operations", r.handlers.GetCustomerOperations)

		// 3.4.4 - Get payment details
		apiRouter.With(standardScheme).Get("/payment/details", r.handlers.GetPaymentDetails)

		// 3.4.5 - Refund payment
		apiRouter.With(standardScheme).Post("/payment/return", r.handlers.RefundPayment)

		enhancedScheme := middleware2.SchemeMiddleware(r.scheme, "enhanced")

		// 4.2.2 - Get trade points (enhanced)
		apiRouter.With(enhancedScheme).Get("/tradepoints/enhanced/{organizationBin}", r.handlers.GetTradePointsEnhanced)

		// 4.2.3 - Register device (enhanced)
		apiRouter.With(enhancedScheme).Post("/device/register/enhanced", r.handlers.RegisterDeviceEnhanced)

		// 4.2.4 - Delete device (enhanced)
		apiRouter.With(enhancedScheme).Post("/device/delete/enhanced", r.handlers.DeleteDeviceEnhanced)

		// 4.3.1 - Create QR code (enhanced)
		apiRouter.With(enhancedScheme).Post("/qr/create/enhanced", r.handlers.CreateQREnhanced)

		// 4.3.2 - Create payment link (enhanced)
		apiRouter.With(enhancedScheme).Post("/qr/create-link/enhanced", r.handlers.CreatePaymentLinkEnhanced)

		// 4.5 - Enhanced refund payment (without customer)
		apiRouter.With(enhancedScheme).Post("/enhanced/payment/return", r.handlers.RefundPaymentEnhanced)

		// 4.6.1 - Get client info by phone number
		apiRouter.With(enhancedScheme).Get("/remote/client-info", r.handlers.GetClientInfo)

		// 4.6.2 - Create remote payment
		apiRouter.With(enhancedScheme).Post("/remote/create", r.handlers.CreateRemotePayment)

		// 4.6.3 - Cancel remote payment
		apiRouter.With(enhancedScheme).Post("/remote/cancel", r.handlers.CancelRemotePayment)

		router.Route("/test", func(apiRouter chi.Router) {
			// 5.1 - Healthcheck
			apiRouter.Get("/health", r.handlers.HealthCheckKaspi)

			// 5.2 - Test QR scan
			apiRouter.Post("/payment/scan", r.handlers.TestScanQR)

			// 5.3 - Test payment confirmation
			apiRouter.Post("/payment/confirm", r.handlers.TestConfirmPayment)

			// 5.4 - Test QR scan error
			apiRouter.Post("/payment/scanerror", r.handlers.TestScanError)

			// 5.5 - Test payment confirmation error
			apiRouter.Post("/payment/confirmerror", r.handlers.TestConfirmError)
		})
	})

	return router
}
