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
)

type MockRefundEnhancedProvider struct {
	RefundPaymentEnhancedFunc func(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error)
	GetClientInfoFunc         func(ctx context.Context, phoneNumber string, deviceToken int64) (*domain.ClientInfoResponse, error)
	CreateRemotePaymentFunc   func(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error)
	CancelRemotePaymentFunc   func(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error)
}

func (m *MockRefundEnhancedProvider) RefundPaymentEnhanced(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error) {
	return m.RefundPaymentEnhancedFunc(ctx, req)
}

func (m *MockRefundEnhancedProvider) GetClientInfo(ctx context.Context, phoneNumber string, deviceToken int64) (*domain.ClientInfoResponse, error) {
	return m.GetClientInfoFunc(ctx, phoneNumber, deviceToken)
}

func (m *MockRefundEnhancedProvider) CreateRemotePayment(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error) {
	return m.CreateRemotePaymentFunc(ctx, req)
}

func (m *MockRefundEnhancedProvider) CancelRemotePayment(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error) {
	return m.CancelRemotePaymentFunc(ctx, req)
}

func TestRefundPaymentEnhancedHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully refunds payment without customer", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			RefundPaymentEnhancedFunc: func(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.QrPaymentID != 123 {
					return nil, fmt.Errorf("invalid payment ID")
				}

				if req.Amount != 10.00 {
					return nil, fmt.Errorf("invalid amount")
				}

				if req.OrganizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return &domain.RefundResponse{
					ReturnOperationID: 20,
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"DeviceToken": "test-token",
			"QrPaymentId": 123,
			"Amount": 10.00,
			"OrganizationBin": "180340021791"
		}`
		req, err := http.NewRequest("POST", "/enhanced/payment/return", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RefundPaymentEnhanced(recorder, req)

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

		var refundResp domain.RefundResponse
		err = json.Unmarshal(jsonData, &refundResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal refund response: %v", err)
		}

		if refundResp.ReturnOperationID != 20 {
			t.Errorf("Expected ReturnOperationID 20, got %d", refundResp.ReturnOperationID)
		}
	})

	t.Run("rejects missing OrganizationBin", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"DeviceToken": "test-token",
			"QrPaymentId": 123,
			"Amount": 10.00
		}`
		req, err := http.NewRequest("POST", "/enhanced/payment/return", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.RefundPaymentEnhanced(recorder, req)

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

func TestGetClientInfoHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets client info", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			GetClientInfoFunc: func(ctx context.Context, phoneNumber string, deviceToken int64) (*domain.ClientInfoResponse, error) {
				if phoneNumber != "87071234567" {
					return nil, fmt.Errorf("invalid phone number")
				}

				if deviceToken != 2 {
					return nil, fmt.Errorf("invalid device token")
				}

				return &domain.ClientInfoResponse{
					ClientName: "John Doe",
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		req, err := http.NewRequest("GET", "/api/remote/client-info?phoneNumber=87071234567&deviceToken=2", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.GetClientInfo(recorder, req)

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

		var clientInfo domain.ClientInfoResponse
		err = json.Unmarshal(jsonData, &clientInfo)
		if err != nil {
			t.Fatalf("Failed to unmarshal client info: %v", err)
		}

		if clientInfo.ClientName != "John Doe" {
			t.Errorf("Expected ClientName John Doe, got %s", clientInfo.ClientName)
		}
	})

	t.Run("rejects missing parameters", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		req, err := http.NewRequest("GET", "/api/remote/client-info?phoneNumber=87071234567", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.GetClientInfo(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}

func TestCreateRemotePaymentHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully creates remote payment", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			CreateRemotePaymentFunc: func(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error) {
				if req.DeviceToken != 2 {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.Amount != 100.00 {
					return nil, fmt.Errorf("invalid amount")
				}

				if req.PhoneNumber != "87071234567" {
					return nil, fmt.Errorf("invalid phone number")
				}

				if req.OrganizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return &domain.RemotePaymentResponse{
					QrPaymentID: 15,
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"OrganizationBin": "180340021791",
			"Amount": 100.00,
			"PhoneNumber": "87071234567",
			"DeviceToken": 2,
			"Comment": "Test payment"
		}`
		req, err := http.NewRequest("POST", "/remote/create", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreateRemotePayment(recorder, req)

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

		var remoteResp domain.RemotePaymentResponse
		err = json.Unmarshal(jsonData, &remoteResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal remote payment response: %v", err)
		}

		if remoteResp.QrPaymentID != 15 {
			t.Errorf("Expected QrPaymentID 15, got %d", remoteResp.QrPaymentID)
		}
	})

	t.Run("rejects missing PhoneNumber", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"OrganizationBin": "180340021791",
			"Amount": 100.00,
			"DeviceToken": 2
		}`
		req, err := http.NewRequest("POST", "/remote/create", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CreateRemotePayment(recorder, req)

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

		expectedError := "PhoneNumber is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestCancelRemotePaymentHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully cancels remote payment", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			CancelRemotePaymentFunc: func(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error) {
				if req.DeviceToken != 2 {
					return nil, fmt.Errorf("invalid device token")
				}

				if req.QrPaymentID != 15 {
					return nil, fmt.Errorf("invalid payment ID")
				}

				if req.OrganizationBin != "180340021791" {
					return nil, fmt.Errorf("invalid organization BIN")
				}

				return &domain.RemotePaymentCancelResponse{
					Status: "RemotePaymentCanceled",
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"OrganizationBin": "180340021791",
			"QrPaymentId": 15,
			"DeviceToken": 2
		}`
		req, err := http.NewRequest("POST", "/remote/cancel", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CancelRemotePayment(recorder, req)

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

		var cancelResp domain.RemotePaymentCancelResponse
		err = json.Unmarshal(jsonData, &cancelResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal cancel response: %v", err)
		}

		if cancelResp.Status != "RemotePaymentCanceled" {
			t.Errorf("Expected Status RemotePaymentCanceled, got %s", cancelResp.Status)
		}
	})

	t.Run("rejects missing QrPaymentId", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"OrganizationBin": "180340021791",
			"DeviceToken": 2
		}`
		req, err := http.NewRequest("POST", "/remote/cancel", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CancelRemotePayment(recorder, req)

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

		expectedError := "QrPaymentId is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})

	t.Run("handles service error", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			CancelRemotePaymentFunc: func(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -99000001,
					Message:    "Purchase with specified identifier not found",
				}
			},
		}

		h := httphandler.NewHandlers(log, nil, nil, nil, nil, nil, nil, mockProvider)

		reqBody := `{
			"OrganizationBin": "180340021791",
			"QrPaymentId": 999,
			"DeviceToken": 2
		}`
		req, err := http.NewRequest("POST", "/remote/cancel", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		h.CancelRemotePayment(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, recorder.Code)
		}

		var resp httphandler.Response
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "Payment with the specified ID not found"
		if !strings.Contains(resp.Error, expectedError) {
			t.Errorf("Expected error message to contain '%s', got '%s'", expectedError, resp.Error)
		}
	})
}
