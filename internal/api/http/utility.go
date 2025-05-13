package http

import (
	"kaspi-api-wrapper/internal/domain"
	"net/http"
)

// HealthCheckKaspi handles health check requests (5.1)
func (h *Handlers) HealthCheckKaspi(w http.ResponseWriter, r *http.Request) {
	err := h.utilityProvider.HealthCheck(r.Context())
	if err != nil {
		h.log.Error("health check failed", "error", err.Error())
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
		h.log.Error("failed to simulate QR scan", "error", err.Error())
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
		h.log.Error("failed to simulate payment confirmation", "error", err.Error())
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
		h.log.Error("failed to simulate QR scan error", "error", err.Error())
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
		h.log.Error("failed to simulate payment confirmation error", "error", err.Error())
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
