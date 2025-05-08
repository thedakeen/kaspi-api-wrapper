package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"strconv"
)

type RefundProvider interface {
	CreateRefundQR(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error)
	GetRefundStatus(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error)
	GetCustomerOperations(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error)
	GetPaymentDetails(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error)
	RefundPayment(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error)
}

func (h *Handlers) CreateRefundQR(w http.ResponseWriter, r *http.Request) {
	var req domain.QRRefundCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
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

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.QrReturnID == 0 {
		BadRequestError(w, "QrReturnId is required")
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

	if qrPaymentIDStr == "" {
		BadRequestError(w, "QrPaymentId is required")
		return
	}

	if deviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	qrPaymentID, err := strconv.ParseInt(qrPaymentIDStr, 10, 64)
	if err != nil {
		BadRequestError(w, "Invalid payment ID format")
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

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.QrPaymentID == 0 {
		BadRequestError(w, "QrPaymentId is required")
		return
	}

	if req.QrReturnID == 0 {
		BadRequestError(w, "QrReturnId is required")
		return
	}

	if req.Amount <= 0 {
		BadRequestError(w, "Amount must be greater than zero")
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
