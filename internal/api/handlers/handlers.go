package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Response is the standard API response format
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Handlers contains all HTTP handlers for the API
type Handlers struct {
	log *slog.Logger
	//kaspiSvc *service.KaspiService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(log *slog.Logger) *Handlers {
	return &Handlers{
		log: log,
		//kaspiSvc: kaspiSvc,
	}
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"failed to marshal JSON response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	resp := Response{
		Success: false,
		Error:   message,
	}
	respondJSON(w, status, resp)
}
