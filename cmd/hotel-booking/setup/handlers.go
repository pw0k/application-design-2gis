package setup

import (
	"application-design-test/internal/handler"
	"github.com/go-chi/chi/v5"
)

func InitHandlers(router *chi.Mux, services *Services) {
	orderHandler := handler.NewOrderHandler(services.BookingService)
	orderHandler.RegisterRoutes(router)

	pingHandler := handler.NewPingHandler()
	pingHandler.RegisterRoutes(router)
}
