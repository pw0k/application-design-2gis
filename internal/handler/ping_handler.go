package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) RegisterRoutes(r chi.Router) {
	r.Get("/ping", h.Ping)
}

func (h *PingHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	message := map[string]string{"message": "pong"}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		return
	}
}
