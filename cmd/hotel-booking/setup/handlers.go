package setup

import (
	"github.com/go-chi/chi/v5"

	"application-design-test/internal/handler"
)

func InitHandlers(router *chi.Mux, services *Services) {
	orderHandler := handler.NewOrderHandler(services.BookingService)
	orderHandler.RegisterRoutes(router)

	pingHandler := handler.NewPingHandler()
	pingHandler.RegisterRoutes(router)
}
