package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	httphandler "kaspi-api-wrapper/internal/api/http"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockPaymentEnhancedProvider struct {
	CreateQREnhancedFunc          func(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error)
	CreatePaymentLinkEnhancedFunc func(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error)
}

func (m *MockPaymentEnhancedProvider) CreateQREnhanced(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error) {
	return m.CreateQREnhancedFunc(ctx, req)
}

func (m *MockPaymentEnhancedProvider) CreatePaymentLinkEnhanced(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
	return m.CreatePaymentLinkEnhancedFunc(ctx, req)
}

func TestCreateQREnhancedHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully creates QR for enhanced scheme", func(t *testing.T) {
		mockProvider := &MockPaymentEnhancedProvider{
			CreateQREnhancedFunc: func(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.Amount != 200.00 {
					return nil, fmt.Errorf("invalid amount")
				}

				if req.OrganizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return &domain.QRCreateResponse{
					QrToken:        "51236903777280167836178166503744993984459",
					ExpireDate:     time.Now().Add(5 * time.Minute),
					QrPaymentID:    15,
					PaymentMethods: []string{"Gold", "Red", "Loan"},
					QrPaymentBehaviorOptions: domain.QRPaymentBehaviorOptions{
						StatusPollingInterval:      5,
						QrCodeScanWaitTimeout:      180,
						PaymentConfirmationTimeout: 65,
					},
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, mockProvider, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"Amount": 200.00,
			"ExternalId": "15",
			"OrganizationBin": "180340021791"
		}`
		req, err := http.NewRequest("POST", "/qr/create/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreateQREnhanced(recorder, req)

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

		var qrResp domain.QRCreateResponse
		err = json.Unmarshal(jsonData, &qrResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal QR response: %v", err)
		}

		if qrResp.QrToken != "51236903777280167836178166503744993984459" {
			t.Errorf("Expected QR token 51236903777280167836178166503744993984459, got %s", qrResp.QrToken)
		}

		if qrResp.QrPaymentID != 15 {
			t.Errorf("Expected QrPaymentID 15, got %d", qrResp.QrPaymentID)
		}
	})

	t.Run("rejects missing OrganizationBin", func(t *testing.T) {
		mockProvider := &MockPaymentEnhancedProvider{}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, mockProvider, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"Amount": 200.00,
			"ExternalId": "15"
		}`
		req, err := http.NewRequest("POST", "/api/qr/create/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreateQREnhanced(recorder, req)

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

		expectedError := "OrganizationBin is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestCreatePaymentLinkEnhancedHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully creates payment link for enhanced scheme", func(t *testing.T) {
		mockProvider := &MockPaymentEnhancedProvider{
			CreatePaymentLinkEnhancedFunc: func(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.Amount != 200.00 {
					return nil, fmt.Errorf("invalid amount")
				}

				if req.OrganizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return &domain.PaymentLinkCreateResponse{
					PaymentLink:    "https://pay.kaspi.kz/pay/123456789",
					ExpireDate:     time.Now().Add(5 * time.Minute),
					PaymentID:      15,
					PaymentMethods: []string{"Gold", "Red", "Loan"},
					PaymentBehaviorOptions: domain.PaymentBehaviorOptions{
						StatusPollingInterval:      5,
						LinkActivationWaitTimeout:  180,
						PaymentConfirmationTimeout: 65,
					},
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, mockProvider, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"Amount": 200.00,
			"ExternalId": "15",
			"OrganizationBin": "180340021791"
		}`
		req, err := http.NewRequest("POST", "/api/qr/create-link/enhanced", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreatePaymentLinkEnhanced(recorder, req)

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

		var linkResp domain.PaymentLinkCreateResponse
		err = json.Unmarshal(jsonData, &linkResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal payment link response: %v", err)
		}

		if linkResp.PaymentLink != "https://pay.kaspi.kz/pay/123456789" {
			t.Errorf("Expected payment link https://pay.kaspi.kz/pay/123456789, got %s", linkResp.PaymentLink)
		}

		if linkResp.PaymentID != 15 {
			t.Errorf("Expected PaymentID 15, got %d", linkResp.PaymentID)
		}
	})
}
