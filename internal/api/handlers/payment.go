package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"strconv"
)

type PaymentProvider interface {
	CreateQR(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error)
	CreatePaymentLink(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error)
	GetPaymentStatus(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error)
}

// CreateQR handles QR code creation for payment (2.3.1)
func (h *Handlers) CreateQR(w http.ResponseWriter, r *http.Request) {
	var req domain.QRCreateRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.Amount <= 0 {
		BadRequestError(w, "Amount must be greater than zero")
		return
	}

	resp, err := h.paymentProvider.CreateQR(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create QR token", err)
		HandleKaspiError(w, err, h.log)
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

	// Validate request
	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.Amount <= 0 {
		BadRequestError(w, "Amount must be greater than zero")
		return
	}

	resp, err := h.paymentProvider.CreatePaymentLink(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create payment link", "error", err)
		HandleKaspiError(w, err, h.log)
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
	if err != nil {
		BadRequestError(w, "Invalid payment ID format")
		return
	}

	status, err := h.paymentProvider.GetPaymentStatus(r.Context(), qrPaymentID)
	if err != nil {
		h.log.Error("failed to get payment status", "error", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    status,
	})
}
