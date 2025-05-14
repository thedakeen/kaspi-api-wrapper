package utility_test

import (
	"context"
	"encoding/json"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	httphandler "kaspi-api-wrapper/internal/handlers/http"
	"kaspi-api-wrapper/internal/validator"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

type MockUtilityProvider struct {
	HealthCheckFunc        func(ctx context.Context) error
	TestScanQRFunc         func(ctx context.Context, req domain.TestScanRequest) error
	TestConfirmPaymentFunc func(ctx context.Context, req domain.TestConfirmRequest) error
	TestScanErrorFunc      func(ctx context.Context, req domain.TestScanErrorRequest) error
	TestConfirmErrorFunc   func(ctx context.Context, req domain.TestConfirmErrorRequest) error
}

func (m *MockUtilityProvider) HealthCheck(ctx context.Context) error {
	return m.HealthCheckFunc(ctx)
}

func (m *MockUtilityProvider) TestScanQR(ctx context.Context, req domain.TestScanRequest) error {
	return m.TestScanQRFunc(ctx, req)
}

func (m *MockUtilityProvider) TestConfirmPayment(ctx context.Context, req domain.TestConfirmRequest) error {
	return m.TestConfirmPaymentFunc(ctx, req)
}

func (m *MockUtilityProvider) TestScanError(ctx context.Context, req domain.TestScanErrorRequest) error {
	return m.TestScanErrorFunc(ctx, req)
}

func (m *MockUtilityProvider) TestConfirmError(ctx context.Context, req domain.TestConfirmErrorRequest) error {
	return m.TestConfirmErrorFunc(ctx, req)
}

func TestHealthCheckKaspi(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully checks health", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			HealthCheckFunc: func(ctx context.Context) error {
				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		req, err := http.NewRequest("GET", "/test/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.HealthCheckKaspi(recorder, req)

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

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map, got %T", resp.Data)
		}

		status, ok := data["status"]
		if !ok {
			t.Fatal("Expected status field in response data")
		}

		if status != "ok" {
			t.Errorf("Expected status 'ok', got '%v'", status)
		}
	})

	t.Run("handles error", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			HealthCheckFunc: func(ctx context.Context) error {
				return &domain.KaspiError{StatusCode: -999, Message: "Service unavailable"}
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		req, err := http.NewRequest("GET", "/test/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.HealthCheckKaspi(recorder, req)

		if recorder.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "Kaspi Pay service is temporarily unavailable"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestTestScanQR(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully simulates QR scan", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			TestScanQRFunc: func(ctx context.Context, req domain.TestScanRequest) error {
				if req.QrPaymentID != "123456" {
					return fmt.Errorf("invalid QR payment ID")
				}
				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		reqBody := `{"qrPaymentId": "123456"}`
		req, err := http.NewRequest("POST", "/test/payment/scan", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.TestScanQR(recorder, req)

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

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map, got %T", resp.Data)
		}

		message, ok := data["message"]
		if !ok {
			t.Fatal("Expected message field in response data")
		}

		expectedMessage := "QR scan simulation successful"
		if message != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
		}
	})

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			TestScanQRFunc: func(ctx context.Context, req domain.TestScanRequest) error {
				return &validator.ValidationError{
					Field:   "qrPaymentId",
					Message: "QR payment ID is required",
					Err:     validator.ErrRequiredField,
				}
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		reqBody := `{"qrPaymentId": ""}`
		req, err := http.NewRequest("POST", "/test/payment/scan", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.TestScanQR(recorder, req)

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

		expectedError := "qrPaymentId: QR payment ID is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestTestConfirmPayment(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully simulates payment confirmation", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			TestConfirmPaymentFunc: func(ctx context.Context, req domain.TestConfirmRequest) error {
				if req.QrPaymentID != "123456" {
					return fmt.Errorf("invalid QR payment ID")
				}
				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		reqBody := `{"qrPaymentId": "123456"}`
		req, err := http.NewRequest("POST", "/test/payment/confirm", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.TestConfirmPayment(recorder, req)

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

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map, got %T", resp.Data)
		}

		expectedMessage := "Payment confirmation simulation successful"
		if data["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, data["message"])
		}
	})
}

func TestTestScanError(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully simulates QR scan error", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			TestScanErrorFunc: func(ctx context.Context, req domain.TestScanErrorRequest) error {
				if req.QrPaymentID != "123456" {
					return fmt.Errorf("invalid QR payment ID")
				}
				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		reqBody := `{"qrPaymentId": "123456"}`
		req, err := http.NewRequest("POST", "/test/payment/scanerror", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.TestScanError(recorder, req)

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

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map, got %T", resp.Data)
		}

		expectedMessage := "QR scan error simulation successful"
		if data["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, data["message"])
		}
	})
}

func TestTestConfirmError(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully simulates payment confirmation error", func(t *testing.T) {
		mockProvider := &MockUtilityProvider{
			TestConfirmErrorFunc: func(ctx context.Context, req domain.TestConfirmErrorRequest) error {
				if req.QrPaymentID != "123456" {
					return fmt.Errorf("invalid QR payment ID")
				}
				return nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, mockProvider, nil, nil, nil, nil)

		reqBody := `{"qrPaymentId": "123456"}`
		req, err := http.NewRequest("POST", "/test/payment/confirmerror", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.TestConfirmError(recorder, req)

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

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map, got %T", resp.Data)
		}

		expectedMessage := "Payment confirmation error simulation successful"
		if data["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, data["message"])
		}
	})
}
