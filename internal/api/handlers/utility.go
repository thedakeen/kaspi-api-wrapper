package handlers

import (
	"context"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
)

type UtilityProvider interface {
	HealthCheck(ctx context.Context) error
	TestScanQR(ctx context.Context, req domain.TestScanRequest) error
	TestConfirmPayment(ctx context.Context, req domain.TestConfirmRequest) error
	TestScanError(ctx context.Context, req domain.TestScanErrorRequest) error
	TestConfirmError(ctx context.Context, req domain.TestConfirmErrorRequest) error
}

// HealthCheckKaspi handles health check requests (5.1)
func (h *Handlers) HealthCheckKaspi(w http.ResponseWriter, r *http.Request) {
	err := h.utilityProvider.HealthCheck(r.Context())
	if err != nil {
		h.log.Error("health check failed", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]string{
			"status": "ok",
		},
	})
}

// TestScanQR handles QR scan simulation (5.2)
func (h *Handlers) TestScanQR(w http.ResponseWriter, r *http.Request) {
	var req domain.TestScanRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.QrPaymentID == "" {
		BadRequestError(w, "qrPaymentId is required")
		return
	}

	err := h.utilityProvider.TestScanQR(r.Context(), req)
	if err != nil {
		h.log.Error("failed to simulate QR scan", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]string{
			"message": "QR scan simulation successful",
		},
	})
}

// TestConfirmPayment handles payment confirmation simulation (5.3)
func (h *Handlers) TestConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var req domain.TestConfirmRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.QrPaymentID == "" {
		BadRequestError(w, "qrPaymentId is required")
		return
	}

	err := h.utilityProvider.TestConfirmPayment(r.Context(), req)
	if err != nil {
		h.log.Error("failed to simulate payment confirmation", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]string{
			"message": "Payment confirmation simulation successful",
		},
	})
}

// TestScanError handles QR scan error simulation (5.4)
func (h *Handlers) TestScanError(w http.ResponseWriter, r *http.Request) {
	var req domain.TestScanErrorRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.QrPaymentID == "" {
		BadRequestError(w, "qrPaymentId is required")
		return
	}

	err := h.utilityProvider.TestScanError(r.Context(), req)
	if err != nil {
		h.log.Error("failed to simulate QR scan error", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]string{
			"message": "QR scan error simulation successful",
		},
	})
}

// TestConfirmError handles payment confirmation error simulation (5.5)
func (h *Handlers) TestConfirmError(w http.ResponseWriter, r *http.Request) {
	var req domain.TestConfirmErrorRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.QrPaymentID == "" {
		BadRequestError(w, "qrPaymentId is required")
		return
	}

	err := h.utilityProvider.TestConfirmError(r.Context(), req)
	if err != nil {
		h.log.Error("failed to simulate payment confirmation error", err)
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]string{
			"message": "Payment confirmation error simulation successful",
		},
	})
}
