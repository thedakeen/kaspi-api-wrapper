package handlers

import (
	"context"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
)

type RefundEnhancedProvider interface {
	RefundPaymentEnhanced(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error)
	GetClientInfo(ctx context.Context, phoneNumber, deviceToken string) (*domain.ClientInfoResponse, error)
	CreateRemotePayment(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error)
	CancelRemotePayment(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error)
}

// RefundPaymentEnhanced handles enhanced payment refund
func (h *Handlers) RefundPaymentEnhanced(w http.ResponseWriter, r *http.Request) {
	var req domain.EnhancedRefundRequest
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

	if req.Amount <= 0 {
		BadRequestError(w, "Amount must be greater than zero")
		return
	}

	if req.OrganizationBin == "" {
		BadRequestError(w, "OrganizationBin is required")
		return
	}

	resp, err := h.refundEnhancedProvider.RefundPaymentEnhanced(r.Context(), req)
	if err != nil {
		h.log.Error("failed to refund payment (enhanced)", err)
		HandleKaspiError(w, err, h.log)
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

	if phoneNumber == "" {
		BadRequestError(w, "phoneNumber is required")
		return
	}

	if deviceToken == "" {
		BadRequestError(w, "deviceToken is required")
		return
	}

	info, err := h.refundEnhancedProvider.GetClientInfo(r.Context(), phoneNumber, deviceToken)
	if err != nil {
		h.log.Error("failed to get client info", err)
		HandleKaspiError(w, err, h.log)
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

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.PhoneNumber == "" {
		BadRequestError(w, "PhoneNumber is required")
		return
	}

	if req.Amount <= 0 {
		BadRequestError(w, "Amount must be greater than zero")
		return
	}

	if req.OrganizationBin == "" {
		BadRequestError(w, "OrganizationBin is required")
		return
	}

	resp, err := h.refundEnhancedProvider.CreateRemotePayment(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create remote payment", err)
		HandleKaspiError(w, err, h.log)
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

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.QrPaymentID == 0 {
		BadRequestError(w, "QrPaymentId is required")
		return
	}

	if req.OrganizationBin == "" {
		BadRequestError(w, "OrganizationBin is required")
		return
	}

	resp, err := h.refundEnhancedProvider.CancelRemotePayment(r.Context(), req)
	if err != nil {
		h.log.Error("failed to cancel remote payment", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}
