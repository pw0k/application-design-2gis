package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) RegisterRoutes(r chi.Router) {
	r.Get("/ping", h.Ping)
}

func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{"message": "pong"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
