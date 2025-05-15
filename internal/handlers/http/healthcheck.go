package http

import (
	"net/http"
	"time"
)

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"status": "ok"},
	})
}
