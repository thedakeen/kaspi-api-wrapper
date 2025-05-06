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
	})

	return router
}
