package handlers

import (
	"encoding/json"
	"io"
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
	log             *slog.Logger
	deviceProvider  DeviceProvider
	paymentProvider PaymentProvider
	utilityProvider UtilityProvider
	refundProvider  RefundProvider

	deviceEnhancedProvider  DeviceEnhancedProvider
	paymentEnhancedProvider PaymentEnhancedProvider
	refundEnhancedProvider  RefundEnhancedProvider
	//kaspiSvc *service.KaspiService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(
	log *slog.Logger,
	deviceProvider DeviceProvider,
	paymentProvider PaymentProvider,
	utilityProvider UtilityProvider,
	refundProvider RefundProvider,

	deviceEnhancedProvider DeviceEnhancedProvider,
	paymentEnhancedProvider PaymentEnhancedProvider,
	refundEnhancedProvider RefundEnhancedProvider,
) *Handlers {
	return &Handlers{
		log:             log,
		deviceProvider:  deviceProvider,
		paymentProvider: paymentProvider,
		utilityProvider: utilityProvider,
		refundProvider:  refundProvider,

		deviceEnhancedProvider:  deviceEnhancedProvider,
		paymentEnhancedProvider: paymentEnhancedProvider,
		refundEnhancedProvider:  refundEnhancedProvider,
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

func BadRequestError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusBadRequest, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusInternalServerError, message)
}

func NotFoundError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusNotFound, message)
}

func ConflictError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusConflict, message)
}

func ForbiddenError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusForbidden, message)
}

func ServiceUnavailableError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusServiceUnavailable, message)
}

func UnauthorizedError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusUnauthorized, message)
}

// DecodeJSONRequest returns Bad Request status in case of invalid data
func DecodeJSONRequest(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "" {
		BadRequestError(w, "Content-Type must be application/json")
		return false
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		var message string

		switch {
		case err == io.EOF:
			message = "Request body is empty"
		case err.Error() == "http: request body too large":
			message = "Request body exceeds size limit"
		default:
			message = "Invalid request format: " + err.Error()
		}

		BadRequestError(w, message)
		return false
	}

	if decoder.More() {
		BadRequestError(w, "Request body must only contain a single JSON object")
		return false
	}

	return true
}
