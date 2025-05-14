package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	httphandler "kaspi-api-wrapper/internal/handlers/http"
	"kaspi-api-wrapper/internal/validator"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type MockDeviceEnhancedProvider struct {
	GetTradePointsEnhancedFunc func(ctx context.Context, organizationBin string) ([]domain.TradePoint, error)
	RegisterDeviceEnhancedFunc func(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDeviceEnhancedFunc   func(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error
}

func (m *MockDeviceEnhancedProvider) GetTradePointsEnhanced(ctx context.Context, organizationBin string) ([]domain.TradePoint, error) {
	return m.GetTradePointsEnhancedFunc(ctx, organizationBin)
}

func (m *MockDeviceEnhancedProvider) RegisterDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
	return m.RegisterDeviceEnhancedFunc(ctx, req)
}

func (m *MockDeviceEnhancedProvider) DeleteDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error {
	return m.DeleteDeviceEnhancedFunc(ctx, req)
}

func TestGetTradePointsEnhancedHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets trade points for organization", func(t *testing.T) {
		mockProvider := &MockDeviceEnhancedProvider{
			GetTradePointsEnhancedFunc: func(ctx context.Context, organizationBin string) ([]domain.TradePoint, error) {
				if organizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return []domain.TradePoint{
					{TradePointID: 1, TradePointName: "Store 1"},
					{TradePointID: 2, TradePointName: "Store 2"},
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, mockProvider, nil, nil)

		r := chi.NewRouter()
		r.Get("/handlers/tradepoints/enhanced/{organizationBin}", h.GetTradePointsEnhanced)

		req, err := http.NewRequest("GET", "/api/tradepoints/enhanced/180340021791", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !resp.Success {
			t.Errorf("Expected success to be true, got false")
		}

		jsonData, err := json.Marshal(resp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var tradePoints []domain.TradePoint
		err = json.Unmarshal(jsonData, &tradePoints)
		if err != nil {
			t.Fatalf("Failed to unmarshal trade points: %v", err)
		}

		if len(tradePoints) != 2 {
			t.Errorf("Expected 2 trade points, got %d", len(tradePoints))
		}

		if tradePoints[0].TradePointID != 1 || tradePoints[0].TradePointName != "Store 1" {
			t.Errorf("Unexpected trade point data: %+v", tradePoints[0])
		}
	})
}

func TestRegisterDeviceEnhancedHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully registers device with organization", func(t *testing.T) {
		mockProvider := &MockDeviceEnhancedProvider{
			RegisterDeviceEnhancedFunc: func(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
				if req.DeviceID != "TEST-DEVICE" {
					return nil, fmt.Errorf("invalid device ID")
				}

				if req.TradePointID != 1 {
					return nil, fmt.Errorf("invalid trade point ID")
				}

				if req.OrganizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return &domain.DeviceRegisterResponse{
					DeviceToken: "2be4cc91-5895-48f8-8bc2-86c7bd419b3b",
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, mockProvider, nil, nil)

		reqBody := `{
			"DeviceId": "TEST-DEVICE",
			"TradePointId": 1,
			"OrganizationBin": "180340021791"
		}`
		req, err := http.NewRequest("POST", "/api/device/register/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RegisterDeviceEnhanced(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !resp.Success {
			t.Errorf("Expected success to be true, got false")
		}

		jsonData, err := json.Marshal(resp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var deviceResp domain.DeviceRegisterResponse
		err = json.Unmarshal(jsonData, &deviceResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal device response: %v", err)
		}

		if deviceResp.DeviceToken != "2be4cc91-5895-48f8-8bc2-86c7bd419b3b" {
			t.Errorf("Expected device token 2be4cc91-5895-48f8-8bc2-86c7bd419b3b, got %s", deviceResp.DeviceToken)
		}
	})

	t.Run("rejects missing OrganizationBin", func(t *testing.T) {
		mockProvider := &MockDeviceEnhancedProvider{
			RegisterDeviceEnhancedFunc: func(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
				if req.OrganizationBin == "" {
					return nil, &validator.ValidationError{
						Field:   "organizationBin",
						Message: "organization BIN is required",
						Err:     validator.ErrRequiredField,
					}
				}
				return &domain.DeviceRegisterResponse{DeviceToken: "test-token"}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, mockProvider, nil, nil)

		reqBody := `{
        "DeviceId": "TEST-DEVICE",
        "TradePointId": 1
    }`
		req, err := http.NewRequest("POST", "/api/device/register/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RegisterDeviceEnhanced(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "organizationBin: organization BIN is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestDeleteDeviceEnhancedHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully deletes device from organization", func(t *testing.T) {
		mockProvider := &MockDeviceEnhancedProvider{
			DeleteDeviceEnhancedFunc: func(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error {
				if req.DeviceToken != "test-token" {
					return fmt.Errorf("invalid device token")
				}

				if req.OrganizationBin != "180340021791" {
					return fmt.Errorf("invalid organization BIN")
				}

				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, mockProvider, nil, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"OrganizationBin": "180340021791"
		}`
		req, err := http.NewRequest("POST", "/api/device/delete/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.DeleteDeviceEnhanced(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !resp.Success {
			t.Errorf("Expected success to be true, got false")
		}
	})

	t.Run("rejects missing DeviceToken", func(t *testing.T) {
		mockProvider := &MockDeviceEnhancedProvider{
			DeleteDeviceEnhancedFunc: func(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error {
				if req.DeviceToken == "" {
					return &validator.ValidationError{
						Field:   "deviceToken",
						Message: "device token is required",
						Err:     validator.ErrRequiredField,
					}
				}
				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, mockProvider, nil, nil)

		reqBody := `{
            "OrganizationBin": "180340021791"
        }`
		req, err := http.NewRequest("POST", "/api/device/delete/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.DeleteDeviceEnhanced(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "deviceToken: device token is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}
