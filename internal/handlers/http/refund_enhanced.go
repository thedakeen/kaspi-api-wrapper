package http

import (
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"strconv"
)

// RefundPaymentEnhanced handles enhanced payment refund
func (h *Handlers) RefundPaymentEnhanced(w http.ResponseWriter, r *http.Request) {
	var req domain.EnhancedRefundRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	resp, err := h.refundEnhancedProvider.RefundPaymentEnhanced(r.Context(), req)
	if err != nil {
		h.log.Error("failed to refund payment (enhanced)", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// GetClientInfo handles getting client information by phone number
func (h *Handlers) GetClientInfo(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.URL.Query().Get("phoneNumber")
	deviceToken := r.URL.Query().Get("deviceToken")

	deviceTokenInt64, err := strconv.ParseInt(deviceToken, 10, 64)
	if err != nil {
		BadRequestError(w, "deviceToken is invalid")
		return
	}

	info, err := h.refundEnhancedProvider.GetClientInfo(r.Context(), phoneNumber, deviceTokenInt64)
	if err != nil {
		h.log.Error("failed to get client info", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    info,
	})
}

// CreateRemotePayment handles creating a remote payment request
func (h *Handlers) CreateRemotePayment(w http.ResponseWriter, r *http.Request) {
	var req domain.RemotePaymentRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	resp, err := h.refundEnhancedProvider.CreateRemotePayment(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create remote payment", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// CancelRemotePayment handles canceling a remote payment request
func (h *Handlers) CancelRemotePayment(w http.ResponseWriter, r *http.Request) {
	var req domain.RemotePaymentCancelRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	resp, err := h.refundEnhancedProvider.CancelRemotePayment(r.Context(), req)
	if err != nil {
		h.log.Error("failed to cancel remote payment", "error", err.Error())
		HandleError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}
