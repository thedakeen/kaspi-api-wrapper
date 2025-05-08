package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"kaspi-api-wrapper/internal/api/handlers"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockRefundProvider struct {
	CreateRefundQRFunc        func(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error)
	GetRefundStatusFunc       func(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error)
	GetCustomerOperationsFunc func(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error)
	GetPaymentDetailsFunc     func(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error)
	RefundPaymentFunc         func(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error)
}

func (m *MockRefundProvider) CreateRefundQR(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error) {
	return m.CreateRefundQRFunc(ctx, req)
}

func (m *MockRefundProvider) GetRefundStatus(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error) {
	return m.GetRefundStatusFunc(ctx, qrReturnID)
}

func (m *MockRefundProvider) GetCustomerOperations(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error) {
	return m.GetCustomerOperationsFunc(ctx, req)
}

func (m *MockRefundProvider) GetPaymentDetails(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error) {
	return m.GetPaymentDetailsFunc(ctx, qrPaymentID, deviceToken)
}

func (m *MockRefundProvider) RefundPayment(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error) {
	return m.RefundPaymentFunc(ctx, req)
}

func TestCreateRefundQRHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully creates refund QR", func(t *testing.T) {
		expireDate, _ := time.Parse(time.RFC3339, "2023-05-16T10:30:00+06:00")

		mockProvider := &MockRefundProvider{
			CreateRefundQRFunc: func(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				return &domain.QRRefundCreateResponse{
					QrToken:    "51236903777280167836178166503744993984459",
					ExpireDate: expireDate,
					QrReturnID: 15,
					QrRefundBehaviorOptions: domain.QRRefundBehaviorOptions{
						QrCodeScanEventPollingInterval: 5,
						QrCodeScanWaitTimeout:          180,
					},
				}, nil
			},
		}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{"DeviceToken": "test-token", "ExternalId": "15"}`
		req, err := http.NewRequest("POST", "/api/return/create", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreateRefundQR(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp handlers.Response
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

		var qrResp domain.QRRefundCreateResponse
		err = json.Unmarshal(jsonData, &qrResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal QR response: %v", err)
		}

		if qrResp.QrToken != "51236903777280167836178166503744993984459" {
			t.Errorf("Expected QR token 51236903777280167836178166503744993984459, got %s", qrResp.QrToken)
		}

		if qrResp.QrReturnID != 15 {
			t.Errorf("Expected QrReturnID 15, got %d", qrResp.QrReturnID)
		}
	})

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockRefundProvider{}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{"ExternalId": "15"}`
		req, err := http.NewRequest("POST", "/api/return/create", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreateRefundQR(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp handlers.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "DeviceToken is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestGetRefundStatusHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets refund status", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			GetRefundStatusFunc: func(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error) {
				if qrReturnID != 15 {
					return nil, fmt.Errorf("invalid return ID")
				}

				return &domain.RefundStatusResponse{
					Status: "QrTokenCreated",
				}, nil
			},
		}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		r := chi.NewRouter()
		r.Get("/return/status/{qrReturnId}", h.GetRefundStatus)

		req, err := http.NewRequest("GET", "/return/status/15", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp handlers.Response
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

		var statusResp domain.RefundStatusResponse
		err = json.Unmarshal(jsonData, &statusResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal status response: %v", err)
		}

		if statusResp.Status != "QrTokenCreated" {
			t.Errorf("Expected Status QrTokenCreated, got %s", statusResp.Status)
		}
	})
}

func TestGetCustomerOperationsHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets customer operations", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			GetCustomerOperationsFunc: func(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error) {
				if req.DeviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.QrReturnID != 15 {
					return nil, fmt.Errorf("invalid return ID")
				}

				return []domain.CustomerOperation{
					{
						QrPaymentID:     900077110,
						TransactionDate: time.Now(),
						Amount:          1.00,
					},
					{
						QrPaymentID:     900077111,
						TransactionDate: time.Now(),
						Amount:          2.00,
					},
				}, nil
			},
		}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{"DeviceToken": "test-token", "QrReturnId": 15, "MaxResult": 10}`
		req, err := http.NewRequest("POST", "/api/return/operations", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.GetCustomerOperations(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp handlers.Response
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

		var operations []domain.CustomerOperation
		err = json.Unmarshal(jsonData, &operations)
		if err != nil {
			t.Fatalf("Failed to unmarshal operations: %v", err)
		}

		if len(operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(operations))
		}

		if operations[0].QrPaymentID != 900077110 {
			t.Errorf("Expected QrPaymentID 900077110, got %d", operations[0].QrPaymentID)
		}

		if operations[1].QrPaymentID != 900077111 {
			t.Errorf("Expected QrPaymentID 900077111, got %d", operations[1].QrPaymentID)
		}
	})
}

func TestGetPaymentDetailsHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets payment details", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			GetPaymentDetailsFunc: func(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error) {
				if qrPaymentID != 123 {
					return nil, fmt.Errorf("invalid payment ID")
				}

				if deviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				return &domain.PaymentDetailsResponse{
					QrPaymentID:           123,
					TotalAmount:           11.00,
					AvailableReturnAmount: 11.00,
					TransactionDate:       time.Now(),
				}, nil
			},
		}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		req, err := http.NewRequest("GET", "/api/payment/details?QrPaymentId=123&DeviceToken=test-token", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.GetPaymentDetails(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp handlers.Response
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

		var details domain.PaymentDetailsResponse
		err = json.Unmarshal(jsonData, &details)
		if err != nil {
			t.Fatalf("Failed to unmarshal details: %v", err)
		}

		if details.QrPaymentID != 123 {
			t.Errorf("Expected QrPaymentID 123, got %d", details.QrPaymentID)
		}

		if details.TotalAmount != 11.00 {
			t.Errorf("Expected TotalAmount 11.00, got %f", details.TotalAmount)
		}

		if details.AvailableReturnAmount != 11.00 {
			t.Errorf("Expected AvailableReturnAmount 11.00, got %f", details.AvailableReturnAmount)
		}
	})

	t.Run("rejects missing parameters", func(t *testing.T) {
		mockProvider := &MockRefundProvider{}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		req, err := http.NewRequest("GET", "/api/payment/details?QrPaymentId=123", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.GetPaymentDetails(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}

func TestRefundPaymentHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully refunds payment", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			RefundPaymentFunc: func(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.QrPaymentID != 123 {
					return nil, fmt.Errorf("invalid payment ID")
				}

				if req.QrReturnID != 15 {
					return nil, fmt.Errorf("invalid return ID")
				}

				if req.Amount != 10.00 {
					return nil, fmt.Errorf("invalid amount")
				}

				return &domain.RefundResponse{
					ReturnOperationID: 20,
				}, nil
			},
		}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"QrPaymentId": 123,
			"QrReturnId": 15,
			"Amount": 10.00
		}`
		req, err := http.NewRequest("POST", "/api/payment/return", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RefundPayment(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		var resp handlers.Response
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

		var refundResp domain.RefundResponse
		err = json.Unmarshal(jsonData, &refundResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal refund response: %v", err)
		}

		if refundResp.ReturnOperationID != 20 {
			t.Errorf("Expected ReturnOperationID 20, got %d", refundResp.ReturnOperationID)
		}
	})

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockRefundProvider{}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{
			"QrPaymentId": 123,
			"QrReturnId": 15,
			"Amount": 10.00
		}`
		req, err := http.NewRequest("POST", "/api/payment/return", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RefundPayment(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp handlers.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "DeviceToken is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})

	t.Run("rejects invalid amount", func(t *testing.T) {
		mockProvider := &MockRefundProvider{}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"QrPaymentId": 123,
			"QrReturnId": 15,
			"Amount": -10.00
		}`
		req, err := http.NewRequest("POST", "/api/payment/return", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RefundPayment(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp handlers.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "Amount must be greater than zero"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})

	t.Run("handles service error", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			RefundPaymentFunc: func(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error) {
				return nil, &domain.KaspiError{StatusCode: -99000005, Message: "Refund amount exceeds purchase amount"}
			},
		}

		h := handlers.NewHandlers(log, nil, nil, nil, mockProvider, nil, nil, nil)

		reqBody := `{
			"DeviceToken": "test-token",
			"QrPaymentId": 123,
			"QrReturnId": 15,
			"Amount": 1000.00
		}`
		req, err := http.NewRequest("POST", "/api/payment/return", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RefundPayment(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp handlers.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "Refund amount cannot exceed the purchase amount"
		if !strings.Contains(resp.Error, expectedError) {
			t.Errorf("Expected error message to contain '%s', got '%s'", expectedError, resp.Error)
		}
	})
}
