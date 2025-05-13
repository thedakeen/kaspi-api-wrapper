package http_test

import (
	"context"
	"encoding/json"
	httphandler "kaspi-api-wrapper/internal/api/http"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDeviceProvider struct {
	GetTradePointsFunc func(ctx context.Context) ([]domain.TradePoint, error)
	RegisterDeviceFunc func(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDeviceFunc   func(ctx context.Context, deviceToken string) error
}

func (m *MockDeviceProvider) GetTradePoints(ctx context.Context) ([]domain.TradePoint, error) {
	return m.GetTradePointsFunc(ctx)
}

func (m *MockDeviceProvider) RegisterDevice(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
	return m.RegisterDeviceFunc(ctx, req)
}

func (m *MockDeviceProvider) DeleteDevice(ctx context.Context, deviceToken string) error {
	return m.DeleteDeviceFunc(ctx, deviceToken)
}

func TestGetTradePoints(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets trade points", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			GetTradePointsFunc: func(ctx context.Context) ([]domain.TradePoint, error) {
				return []domain.TradePoint{
					{TradePointID: 1, TradePointName: "Store 1"},
					{TradePointID: 2, TradePointName: "Store 2"},
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, mockProvider, nil, nil, nil, nil, nil, nil)

		req, err := createRequest(http.MethodGet, "/api/tradepoints", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.GetTradePoints(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
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

	t.Run("handles error", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			GetTradePointsFunc: func(ctx context.Context) ([]domain.TradePoint, error) {
				return nil, &domain.KaspiError{StatusCode: -14000002, Message: "No trade points available"}
			},
		}

		h := httphandler.NewHandlers(log, mockProvider, nil, nil, nil, nil, nil, nil)

		req, err := createRequest(http.MethodGet, "/api/tradepoints", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.GetTradePoints(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "No trade points available. Please create a trade point in the Kaspi Pay application"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestRegisterDevice(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully registers device", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			RegisterDeviceFunc: func(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
				return &domain.DeviceRegisterResponse{
					DeviceToken: "2be4cc91-5895-48f8-8bc2-86c7bd419b3b",
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, mockProvider, nil, nil, nil, nil, nil, nil)

		registerReq := domain.DeviceRegisterRequest{
			DeviceID:     "TEST-DEVICE",
			TradePointID: 1,
		}

		req, err := createRequest(http.MethodPost, "/api/device/register", registerReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.RegisterDevice(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
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

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{}

		h := httphandler.NewHandlers(log, mockProvider, nil, nil, nil, nil, nil, nil)

		registerReq := domain.DeviceRegisterRequest{
			DeviceID: "TEST-DEVICE",
		}

		req, err := createRequest(http.MethodPost, "/api/device/register", registerReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.RegisterDevice(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "TradePointID is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestDeleteDevice(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully deletes device", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			DeleteDeviceFunc: func(ctx context.Context, deviceToken string) error {
				return nil
			},
		}

		h := httphandler.NewHandlers(log, mockProvider, nil, nil, nil, nil, nil, nil)

		deleteReq := struct {
			DeviceToken string `json:"deviceToken"`
		}{
			DeviceToken: "test-token",
		}

		req, err := createRequest(http.MethodPost, "/api/device/delete", deleteReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.DeleteDevice(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !resp.Success {
			t.Errorf("Expected success to be true, got false")
		}
	})

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{}

		h := httphandler.NewHandlers(log, mockProvider, nil, nil, nil, nil, nil, nil)

		deleteReq := struct {
			DeviceToken string `json:"deviceToken"`
		}{
			DeviceToken: "",
		}

		req, err := createRequest(http.MethodPost, "/api/device/delete", deleteReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.DeleteDevice(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "device token is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}
