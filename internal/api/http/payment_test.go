package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	httphandler "kaspi-api-wrapper/internal/api/http"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/validator"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockPaymentProvider struct {
	CreateQRFunc          func(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error)
	CreatePaymentLinkFunc func(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error)
	GetPaymentStatusFunc  func(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error)
}

func (m *MockPaymentProvider) CreateQR(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error) {
	if m.CreateQRFunc != nil {
		return m.CreateQRFunc(ctx, req)
	}

	if req.DeviceToken == "" {
		return nil, &validator.ValidationError{
			Field:   "deviceToken",
			Message: "device token is required",
			Err:     validator.ErrRequiredField,
		}
	}

	if req.Amount <= 0 {
		return nil, &validator.ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     validator.ErrInvalidAmount,
		}
	}

	return &domain.QRCreateResponse{
		QrToken:     "test-qr-token",
		QrPaymentID: 123,
	}, nil
}

func (m *MockPaymentProvider) CreatePaymentLink(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
	if m.CreatePaymentLinkFunc != nil {
		return m.CreatePaymentLinkFunc(ctx, req)
	}

	if req.DeviceToken == "" {
		return nil, &validator.ValidationError{
			Field:   "deviceToken",
			Message: "device token is required",
			Err:     validator.ErrRequiredField,
		}
	}

	if req.Amount <= 0 {
		return nil, &validator.ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     validator.ErrInvalidAmount,
		}
	}

	expireDate := time.Now().Add(5 * time.Minute)
	return &domain.PaymentLinkCreateResponse{
		PaymentLink:    "https://pay.kaspi.kz/pay/123456789",
		ExpireDate:     expireDate,
		PaymentID:      15,
		PaymentMethods: []string{"Gold", "Red", "Loan"},
		PaymentBehaviorOptions: domain.PaymentBehaviorOptions{
			StatusPollingInterval:      5,
			LinkActivationWaitTimeout:  180,
			PaymentConfirmationTimeout: 65,
		},
	}, nil
}

func (m *MockPaymentProvider) GetPaymentStatus(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error) {
	return m.GetPaymentStatusFunc(ctx, qrPaymentID)
}

func TestCreateQRHandler(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully creates QR", func(t *testing.T) {
		expireDate, _ := time.Parse(time.RFC3339, "2023-05-16T10:30:00+06:00")

		mockProvider := &MockPaymentProvider{
			CreateQRFunc: func(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, errors.New("invalid device token")
				}

				if req.Amount != 200.00 {
					return nil, errors.New("invalid amount")
				}

				return &domain.QRCreateResponse{
					QrToken:        "51236903777280167836178166503744993984459",
					ExpireDate:     expireDate,
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

		h := httphandler.NewHandlers(log, nil, mockProvider, nil, nil, nil, nil, nil)

		createReq := domain.QRCreateRequest{
			DeviceToken: "test-token",
			Amount:      200.00,
			ExternalID:  "15",
		}

		req, err := createRequest(http.MethodPost, "/api/qr/create", createReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.CreateQR(recorder, req)

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

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{}

		h := httphandler.NewHandlers(log, nil, mockProvider, nil, nil, nil, nil, nil)

		createReq := domain.QRCreateRequest{
			DeviceToken: "test-token",
			Amount:      -10.00,
			ExternalID:  "15",
		}

		req, err := createRequest(http.MethodPost, "/api/qr/create", createReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.CreateQR(recorder, req)

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

		expectedError := "amount: amount must be greater than zero"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})
}

func TestCreatePaymentLinkFunc(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully creates payment link", func(t *testing.T) {
		expireDate, _ := time.Parse(time.RFC3339, "2023-05-16T10:30:00+06:00")

		mockProvider := &MockPaymentProvider{
			CreatePaymentLinkFunc: func(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, errors.New("invalid device token")
				}

				if req.Amount != 200.00 {
					return nil, errors.New("invalid amount")
				}

				return &domain.PaymentLinkCreateResponse{
					PaymentLink:    "https://pay.kaspi.kz/pay/123456789",
					ExpireDate:     expireDate,
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

		h := httphandler.NewHandlers(log, nil, mockProvider, nil, nil, nil, nil, nil)

		createReq := domain.PaymentLinkCreateRequest{
			DeviceToken: "test-token",
			Amount:      200.00,
			ExternalID:  "15",
		}

		req, err := createRequest(http.MethodPost, "/api/qr/create-link", createReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.CreatePaymentLink(recorder, req)

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

		if !linkResp.ExpireDate.Equal(expireDate) {
			t.Errorf("Expected ExpireDate %v, got %v", expireDate, linkResp.ExpireDate)
		}

		if len(linkResp.PaymentMethods) != 3 {
			t.Errorf("Expected 3 payment methods, got %d", len(linkResp.PaymentMethods))
		}

		if linkResp.PaymentBehaviorOptions.StatusPollingInterval != 5 {
			t.Errorf("Expected StatusPollingInterval 5, got %d", linkResp.PaymentBehaviorOptions.StatusPollingInterval)
		}
	})

	t.Run("rejects invalid request", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{}

		h := httphandler.NewHandlers(log, nil, mockProvider, nil, nil, nil, nil, nil)

		createReq := domain.PaymentLinkCreateRequest{
			DeviceToken: "",
			Amount:      200.00,
			ExternalID:  "15",
		}

		req, err := createRequest(http.MethodPost, "/api/qr/create-link", createReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.CreatePaymentLink(recorder, req)

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

		expectedError := "deviceToken: device token is required"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})

	t.Run("handles service error", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{
			CreatePaymentLinkFunc: func(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
				return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
			},
		}

		h := httphandler.NewHandlers(log, nil, mockProvider, nil, nil, nil, nil, nil)

		createReq := domain.PaymentLinkCreateRequest{
			DeviceToken: "invalid-token",
			Amount:      200.00,
			ExternalID:  "15",
		}

		req, err := createRequest(http.MethodPost, "/api/qr/create-link", createReq)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		h.CreatePaymentLink(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "Device not found"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})

}

func TestGetPaymentStatus(t *testing.T) {
	log := setupTestLogger()

	t.Run("successfully gets payment status", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{
			GetPaymentStatusFunc: func(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error) {
				if qrPaymentID != 15 {
					return nil, errors.New("invalid payment ID")
				}

				return &domain.PaymentStatusResponse{
					Status:        "Wait",
					TransactionID: "35134863",
					LoanOfferName: "Рассрочка 0-0-12",
					LoanTerm:      12,
					IsOffer:       true,
					ProductType:   "Loan",
					Amount:        200.00,
					StoreName:     "Store 1",
					Address:       "Test Address",
					City:          "Almaty",
				}, nil
			},
		}

		h := httphandler.NewHandlers(log, nil, mockProvider, nil, nil, nil, nil, nil)

		r := chi.NewRouter()
		r.Get("/payment/status/{qrPaymentId}", h.GetPaymentStatus)

		req, err := http.NewRequest("GET", "/payment/status/15", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

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

		var statusResp domain.PaymentStatusResponse
		err = json.Unmarshal(jsonData, &statusResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal status response: %v", err)
		}

		if statusResp.Status != "Wait" {
			t.Errorf("Expected Status Wait, got %s", statusResp.Status)
		}

		if statusResp.TransactionID != "35134863" {
			t.Errorf("Expected TransactionID 35134863, got %s", statusResp.TransactionID)
		}
	})
}
