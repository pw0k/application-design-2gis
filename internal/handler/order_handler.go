package handler

import (
	"application-design-test/internal/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"application-design-test/internal/model"
)

type BookingService interface {
	CreateOrder(ctx context.Context, order *model.Order) error
}

type OrderHandler struct {
	BookingService BookingService
}

func NewOrderHandler(bookingService BookingService) *OrderHandler {
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
	ctx := r.Context()
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Error("Incorrect request", "error", err)
		http.Error(w, "Incorrect request", http.StatusBadRequest)
		return
	}

	if err := validate(order); err != nil {
		slog.Error("Incorrect request", "error", err)
		http.Error(w, "Incorrect request", http.StatusBadRequest)
		return
	}

	if err := h.BookingService.CreateOrder(ctx, &order); err != nil {
		switch {
		case errors.Is(err, service.ErrQuotaUnavailable):
			http.Error(w, "Quota unavailable error", http.StatusConflict)
			slog.Error("Quota unavailable ", "error", err)
		default:
			http.Error(w, "CreateOrder error", http.StatusInternalServerError)
			slog.Error("CreateOrder error", "error", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(order)
	if err != nil {
		return
	}
}

func validate(order model.Order) error {
	if order.From.After(order.To) {
		return fmt.Errorf("incorrect request, order.from is after order.to")
	}
	return nil
}
