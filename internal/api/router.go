package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"kaspi-api-wrapper/internal/api/handlers"
	custommw "kaspi-api-wrapper/internal/api/middleware"
	"log/slog"
)

type Router struct {
	handlers *handlers.Handlers
}

func NewRouter(log *slog.Logger, handlers *handlers.Handlers) *Router {
	return &Router{
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

	return router
}
