package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
)

type DeviceEnhancedProvider interface {
	GetTradePointsEnhanced(ctx context.Context, organizationBin string) ([]domain.TradePoint, error)
	RegisterDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error
}

// GetTradePointsEnhanced handles a request to get trade points in the enhanced scheme (4.2.2)
func (h *Handlers) GetTradePointsEnhanced(w http.ResponseWriter, r *http.Request) {
	organizationBin := chi.URLParam(r, "organizationBin")
	if organizationBin == "" {
		BadRequestError(w, "OrganizationBin is required")
		return
	}

	tradePoints, err := h.deviceEnhancedProvider.GetTradePointsEnhanced(r.Context(), organizationBin)
	if err != nil {
		h.log.Error("failed to get trade points (enhanced)", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tradePoints,
	})
}

// RegisterDeviceEnhanced processes a request to register a device in the enhanced scheme (4.2.3)
func (h *Handlers) RegisterDeviceEnhanced(w http.ResponseWriter, r *http.Request) {
	var req domain.EnhancedDeviceRegisterRequest
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

	if req.OrganizationBin == "" {
		BadRequestError(w, "OrganizationBin is required")
		return
	}

	resp, err := h.deviceEnhancedProvider.RegisterDeviceEnhanced(r.Context(), req)
	if err != nil {
		h.log.Error("failed to register device (enhanced)", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    resp,
	})
}

// DeleteDeviceEnhanced handles a request to delete a device in the enhanced scheme (4.2.4)
func (h *Handlers) DeleteDeviceEnhanced(w http.ResponseWriter, r *http.Request) {
	var req domain.EnhancedDeviceDeleteRequest
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

	if req.DeviceToken == "" {
		BadRequestError(w, "DeviceToken is required")
		return
	}

	if req.OrganizationBin == "" {
		BadRequestError(w, "OrganizationBin is required")
		return
	}

	err := h.deviceEnhancedProvider.DeleteDeviceEnhanced(r.Context(), req)
	if err != nil {
		h.log.Error("failed to delete device (enhanced)", "error", err.Error())
		HandleKaspiError(w, err, h.log)
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"message": "Device deleted successfully"},
	})
}
