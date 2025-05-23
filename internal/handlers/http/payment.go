package http

import (
	"github.com/go-chi/chi/v5"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"strconv"
)

// CreateQR handles QR code creation for payment (2.3.1)
func (h *Handlers) CreateQR(w http.ResponseWriter, r *http.Request) {
	var req domain.QRCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	resp, err := h.paymentProvider.CreateQR(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create QR token", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// CreatePaymentLink handles payment link creation (2.3.2)
func (h *Handlers) CreatePaymentLink(w http.ResponseWriter, r *http.Request) {
	var req domain.PaymentLinkCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	resp, err := h.paymentProvider.CreatePaymentLink(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create payment link", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// GetPaymentStatus handles payment status retrieval (2.3.3)
func (h *Handlers) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	qrPaymentIDStr := chi.URLParam(r, "qrPaymentId")
	qrPaymentID, err := strconv.ParseInt(qrPaymentIDStr, 10, 64)

	status, err := h.paymentProvider.GetPaymentStatus(r.Context(), qrPaymentID)
	if err != nil {
		h.log.Error("failed to get payment status", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    status,
	})
}
