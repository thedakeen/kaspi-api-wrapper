package http

import (
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/validator"
	"net/http"
)

// CreateQREnhanced handles a request to create a QR in the enhanced scheme (4.3.1)
func (h *Handlers) CreateQREnhanced(w http.ResponseWriter, r *http.Request) {
	var req domain.EnhancedQRCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if err := validator.ValidateEnhancedQRCreateRequest(req); err != nil {
		h.log.Warn("invalid enhanced QR create request", "error", err.Error())
		validator.HTTPError(w, err)
		return
	}

	resp, err := h.paymentEnhancedProvider.CreateQREnhanced(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create QR (enhanced)", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// CreatePaymentLinkEnhanced handles a request to create a payment link in the enhanced scheme (4.3.2)
func (h *Handlers) CreatePaymentLinkEnhanced(w http.ResponseWriter, r *http.Request) {
	var req domain.EnhancedPaymentLinkCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if err := validator.ValidateEnhancedPaymentLinkCreateRequest(req); err != nil {
		h.log.Warn("invalid enhanced payment link create request", "error", err.Error())
		validator.HTTPError(w, err)
		return
	}

	resp, err := h.paymentEnhancedProvider.CreatePaymentLinkEnhanced(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create payment link (enhanced)", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}
