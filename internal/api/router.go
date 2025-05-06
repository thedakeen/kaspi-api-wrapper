package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"kaspi-api-wrapper/internal/api/handlers"
	custommw "kaspi-api-wrapper/internal/api/middleware"
	"log/slog"
)

type Router struct {
	log      *slog.Logger
	handlers *handlers.Handlers
}

func NewRouter(log *slog.Logger, handlers *handlers.Handlers) *Router {
	return &Router{
		log:      log,
		handlers: handlers,
	}
}

func (r *Router) Setup() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(custommw.Logger)
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

		router.Route("/test", func(apiRouter chi.Router) {
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
