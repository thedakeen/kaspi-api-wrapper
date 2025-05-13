package http

import (
	"github.com/go-chi/chi/v5"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/validator"
	"net/http"
	"strconv"
)

func (h *Handlers) CreateRefundQR(w http.ResponseWriter, r *http.Request) {
	var req domain.QRRefundCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if err := validator.ValidateQRRefundCreateRequest(req); err != nil {
		h.log.Warn("invalid refund QR create request", "error", err.Error())
		validator.HTTPError(w, err)
		return
	}

	resp, err := h.refundProvider.CreateRefundQR(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create refund QR token", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// GetRefundStatus handles refund status retrieval
func (h *Handlers) GetRefundStatus(w http.ResponseWriter, r *http.Request) {
	qrReturnIDStr := chi.URLParam(r, "qrReturnId")
	qrReturnID, err := strconv.ParseInt(qrReturnIDStr, 10, 64)
	if err != nil {
		BadRequestError(w, "Invalid refund ID format")
		return
	}

	status, err := h.refundProvider.GetRefundStatus(r.Context(), qrReturnID)
	if err != nil {
		h.log.Error("failed to get refund status", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    status,
	})
}

// GetCustomerOperations handles getting customer operations
func (h *Handlers) GetCustomerOperations(w http.ResponseWriter, r *http.Request) {
	var req domain.CustomerOperationsRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if err := validator.ValidateCustomerOperationsRequest(req); err != nil {
		h.log.Warn("invalid customer operations request", "error", err.Error())
		validator.HTTPError(w, err)
		return
	}

	operations, err := h.refundProvider.GetCustomerOperations(r.Context(), req)
	if err != nil {
		h.log.Error("failed to get customer operations", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    operations,
	})
}

// GetPaymentDetails handles getting payment details
func (h *Handlers) GetPaymentDetails(w http.ResponseWriter, r *http.Request) {
	qrPaymentIDStr := r.URL.Query().Get("QrPaymentId")
	deviceToken := r.URL.Query().Get("DeviceToken")

	qrPaymentID, err := strconv.ParseInt(qrPaymentIDStr, 10, 64)
	if err != nil {
		BadRequestError(w, "Invalid payment ID format")
		return
	}

	if err := validator.ValidatePaymentDetailsRequest(qrPaymentID, deviceToken); err != nil {
		h.log.Warn("invalid payment details request", "error", err.Error())
		validator.HTTPError(w, err)
		return
	}

	details, err := h.refundProvider.GetPaymentDetails(r.Context(), qrPaymentID, deviceToken)
	if err != nil {
		h.log.Error("failed to get payment details", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    details,
	})
}

// RefundPayment handles payment refund
func (h *Handlers) RefundPayment(w http.ResponseWriter, r *http.Request) {
	var req domain.RefundRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if err := validator.ValidateRefundRequest(req); err != nil {
		h.log.Warn("invalid refund request", "error", err.Error())
		validator.HTTPError(w, err)
		return
	}

	resp, err := h.refundProvider.RefundPayment(r.Context(), req)
	if err != nil {
		h.log.Error("failed to refund payment", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}
