package handlers

import (
	"context"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
)

type DeviceProvider interface {
	GetTradePoints(ctx context.Context) ([]domain.TradePoint, error)
	RegisterDevice(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDevice(ctx context.Context, deviceToken string) error
}

// AllTradePoints handles retrieving trade points (2.2.2)
func (h *Handlers) AllTradePoints(w http.ResponseWriter, r *http.Request) {
	tradePoints, err := h.deviceProvider.GetTradePoints(r.Context())
	if err != nil {
		h.log.Error("failed to get trade points", err)
		InternalServerError(w, "Failed to get trade points")
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tradePoints,
	})
}

// RegisterDevice handles device registration (2.2.3)
func (h *Handlers) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req domain.DeviceRegisterRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.DeviceID == "" {
		BadRequestError(w, "DeviceID is required")
		return
	}

	if req.TradePointID == 0 {
		BadRequestError(w, "TradePointID is required")
		return
	}

	resp, err := h.deviceProvider.RegisterDevice(r.Context(), req)
	if err != nil {
		h.log.Error("failed to register device", err)
		InternalServerError(w, "failed to register device")
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// DeleteDevice hadnles device deletion (2.2.4)
func (h *Handlers) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceToken string `json:"deviceToken"`
	}

	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.DeviceToken == "" {
		BadRequestError(w, "device token is required")
		return
	}

	err := h.deviceProvider.DeleteDevice(r.Context(), req.DeviceToken)
	if err != nil {
		h.log.Error("failed to delete device", err)
		InternalServerError(w, "failed to delete device")
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"message": "Device deleted successfully"},
	})
}
