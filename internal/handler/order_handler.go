package handler

import (
	//"application-design-test/internal/handler/dto"
	"application-design-test/internal/model"
	"application-design-test/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	BookingService service.BookingService
}

func NewOrderHandler(bookingService service.BookingService) *OrderHandler {
	return &OrderHandler{
		BookingService: bookingService,
	}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", h.CreateOrder)
	})
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Error("Некорректный запрос", "ошибка", err)
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	if err := h.BookingService.CreateOrder(&order); err != nil {
		slog.Error("Не удалось оформить", "ошибка", err)
		//todo: здесь по хорошему надо выдавать разные коды ошибок
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(order)
	if err != nil {
		return
	}
}
