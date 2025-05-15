package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/service"
	"kaspi-api-wrapper/internal/storage"
	"kaspi-api-wrapper/internal/testutils"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func setupTestService(log *slog.Logger, scheme string) (*service.KaspiService, *testutils.MockHTTPClient) {
	mockClient := &testutils.MockHTTPClient{}

	mockSaver := &MockDeviceSaver{
		SaveDeviceFunc: func(ctx context.Context, deviceID, deviceToken string, tradePointID int64) error {
			return nil
		},
		SaveDeviceEnhancedFunc: func(ctx context.Context, deviceID, deviceToken string, tradePointID int64, organizationBin string) error {
			return nil
		},
	}

	svc := service.NewKaspiService(
		log,
		scheme,
		"https://test.com",
		"https://test.com",
		"https://test.com",
		"test-handlers-key",
		nil,
		mockSaver,
	)

	svc.SetHTTPClient(mockClient)

	return svc, mockClient
}

type MockDeviceSaver struct {
	SaveDeviceFunc         func(ctx context.Context, deviceID, deviceToken string, tradePointID int64) error
	SaveDeviceEnhancedFunc func(ctx context.Context, deviceID, deviceToken string, tradePointID int64, organizationBin string) error
}

func (m *MockDeviceSaver) SaveDevice(ctx context.Context, deviceID, deviceToken string, tradePointID int64) error {
	if m.SaveDeviceFunc != nil {
		return m.SaveDeviceFunc(ctx, deviceID, deviceToken, tradePointID)
	}
	return nil
}

func (m *MockDeviceSaver) SaveDeviceEnhanced(ctx context.Context, deviceID, deviceToken string, tradePointID int64, organizationBin string) error {
	if m.SaveDeviceEnhancedFunc != nil {
		return m.SaveDeviceEnhancedFunc(ctx, deviceID, deviceToken, tradePointID, organizationBin)
	}
	return nil
}

func TestGetBaseURL(t *testing.T) {
	t.Run("returns basic URL for basic scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc, _ := setupTestService(log, "basic")

		baseURL := svc.GetBaseURL()
		if baseURL != "https://test.com" {
			t.Errorf("Expected base URL https://test.com, got %s", baseURL)
		}
	})

	t.Run("returns standard URL for standard scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc, _ := setupTestService(log, "standard")

		baseURL := svc.GetBaseURL()
		if baseURL != "https://test.com" {
			t.Errorf("Expected base URL https://test.com, got %s", baseURL)
		}
	})

	t.Run("returns enhanced URL for enhanced scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc, _ := setupTestService(log, "enhanced")

		baseURL := svc.GetBaseURL()
		if baseURL != "https://test.com" {
			t.Errorf("Expected base URL https://test.com, got %s", baseURL)
		}
	})

	t.Run("defaults to basic URL for unknown scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc, _ := setupTestService(log, "unknown")

		baseURL := svc.GetBaseURL()
		if baseURL != "https://test.com" {
			t.Errorf("Expected base URL https://test.com, got %s", baseURL)
		}
	})
}

func TestRequest(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			// Verify request properties
			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			if req.URL.String() != "https://test.com/test-path" {
				t.Errorf("Expected URL https://test.com/test-path, got %s", req.URL.String())
			}

			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type header application/json, got %s",
					req.Header.Get("Content-Type"))
			}

			if req.Header.Get("X-Request-ID") == "" {
				t.Error("Expected X-Request-ID header to be set")
			}

			if req.Header.Get("Api-Key") != "test-handlers-key" {
				t.Errorf("Expected Api-Key header test-handlers-key, got %s",
					req.Header.Get("Api-Key"))
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {"test": "value"}
			}`), nil
		}

		var result map[string]interface{}
		err := svc.Request(context.Background(), http.MethodGet, "/test-path", nil, &result)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result["test"] != "value" {
			t.Errorf("Expected result.test to be 'value', got %v", result["test"])
		}
	})

	t.Run("HTTP client error", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		}

		var result map[string]interface{}
		err := svc.Request(context.Background(), http.MethodGet, "/test-path", nil, &result)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "network error") {
			t.Errorf("Expected network error, got %v", err)
		}
	})

	t.Run("Kaspi API error", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": -1501,
				"Message": "Device not found"
			}`), nil
		}

		var result map[string]interface{}
		err := svc.Request(context.Background(), http.MethodGet, "/test-path", nil, &result)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		kaspiErr, ok := domain.IsKaspiError(err)
		if !ok {
			t.Fatalf("Expected KaspiError, got %T", err)
		}

		if kaspiErr.StatusCode != -1501 {
			t.Errorf("Expected status code -1501, got %d", kaspiErr.StatusCode)
		}

		if kaspiErr.Message != "Device not found" {
			t.Errorf("Expected message 'Device not found', got '%s'", kaspiErr.Message)
		}
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return testutils.NewMockResponse(http.StatusOK, `invalid json`), nil
		}

		var result map[string]interface{}
		err := svc.Request(context.Background(), http.MethodGet, "/test-path", nil, &result)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "unmarshal") {
			t.Errorf("Expected JSON unmarshal error, got %v", err)
		}
	})
}

//////// 	Device operations testing		////////

func TestGetTradePoints(t *testing.T) {
	t.Run("successfully gets trade points", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/partner/tradepoints" {
				t.Errorf("Expected URL path /partner/tradepoints, got %s", req.URL.Path)
			}

			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": [
					{"TradePointId": 1, "TradePointName": "Store 1"},
					{"TradePointId": 2, "TradePointName": "Store 2"}
				]
			}`), nil
		}

		tradePoints, err := svc.GetTradePoints(context.Background())

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(tradePoints) != 2 {
			t.Errorf("Expected 2 trade points, got %d", len(tradePoints))
		}

		if tradePoints[0].TradePointID != 1 || tradePoints[0].TradePointName != "Store 1" {
			t.Errorf("Unexpected trade point data: %+v", tradePoints[0])
		}

		if tradePoints[1].TradePointID != 2 || tradePoints[1].TradePointName != "Store 2" {
			t.Errorf("Unexpected trade point data: %+v", tradePoints[1])
		}
	})

	t.Run("handles error response", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": -14000002,
				"Message": "No trade points available"
			}`), nil
		}

		_, err := svc.GetTradePoints(context.Background())

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		kaspiErr, ok := domain.IsKaspiError(err)
		if !ok {
			t.Fatalf("Expected KaspiError, got %T", err)
		}

		if kaspiErr.StatusCode != -14000002 {
			t.Errorf("Expected status code -14000002, got %d", kaspiErr.StatusCode)
		}
	})
}

func TestRegisterDevice(t *testing.T) {
	t.Run("successfully registers device", func(t *testing.T) {
		log := setupTestLogger()

		var savedDeviceID string
		var savedDeviceToken string
		var savedTradePointID int64
		saveDeviceCalled := false

		mockSaver := &MockDeviceSaver{
			SaveDeviceFunc: func(ctx context.Context, deviceID, deviceToken string, tradePointID int64) error {
				saveDeviceCalled = true
				savedDeviceID = deviceID
				savedDeviceToken = deviceToken
				savedTradePointID = tradePointID
				return nil
			},
		}

		mockClient := &testutils.MockHTTPClient{}

		svc := service.NewKaspiService(
			log,
			"basic",
			"https://test.com",
			"https://test.com",
			"https://test.com",
			"test-handlers-key",
			nil,
			mockSaver,
		)

		svc.SetHTTPClient(mockClient)

		expectedDeviceToken := "2be4cc91-5895-48f8-8bc2-86c7bd419b3b"
		expectedDeviceID := "TEST-DEVICE"
		expectedTradePointID := int64(1)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/device/register" {
				t.Errorf("Expected URL path /device/register, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			body, _ := io.ReadAll(req.Body)
			req.Body.Close()

			var reqBody domain.DeviceRegisterRequest
			err := json.Unmarshal(body, &reqBody)
			if err != nil {
				t.Errorf("Failed to parse request body: %v", err)
			}

			if reqBody.DeviceID != expectedDeviceID {
				t.Errorf("Expected DeviceID %s, got %s", expectedDeviceID, reqBody.DeviceID)
			}

			if reqBody.TradePointID != expectedTradePointID {
				t.Errorf("Expected TradePointID %d, got %d", expectedTradePointID, reqBody.TradePointID)
			}

			return testutils.NewMockResponse(http.StatusOK, fmt.Sprintf(`{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"DeviceToken": "%s"
				}
			}`, expectedDeviceToken)), nil
		}

		resp, err := svc.RegisterDevice(context.Background(), domain.DeviceRegisterRequest{
			DeviceID:     expectedDeviceID,
			TradePointID: expectedTradePointID,
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.DeviceToken != expectedDeviceToken {
			t.Errorf("Expected device token %s, got %s", expectedDeviceToken, resp.DeviceToken)
		}

		if !saveDeviceCalled {
			t.Error("SaveDevice method was not called")
		}

		if savedDeviceID != expectedDeviceID {
			t.Errorf("SaveDevice called with wrong deviceID, expected: %s, got: %s",
				expectedDeviceID, savedDeviceID)
		}

		if savedDeviceToken != expectedDeviceToken {
			t.Errorf("SaveDevice called with wrong deviceToken, expected: %s, got: %s",
				expectedDeviceToken, savedDeviceToken)
		}

		if savedTradePointID != expectedTradePointID {
			t.Errorf("SaveDevice called with wrong tradePointID, expected: %d, got: %d",
				expectedTradePointID, savedTradePointID)
		}
	})

	t.Run("handles device already exists error", func(t *testing.T) {
		log := setupTestLogger()

		mockSaver := &MockDeviceSaver{
			SaveDeviceFunc: func(ctx context.Context, deviceID, deviceToken string, tradePointID int64) error {
				return storage.ErrDeviceExists
			},
		}

		mockClient := &testutils.MockHTTPClient{}

		svc := service.NewKaspiService(
			log,
			"basic",
			"https://test.com",
			"https://test.com",
			"https://test.com",
			"test-handlers-key",
			nil,
			mockSaver,
		)

		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"DeviceToken": "2be4cc91-5895-48f8-8bc2-86c7bd419b3b"
				}
			}`), nil
		}

		_, err := svc.RegisterDevice(context.Background(), domain.DeviceRegisterRequest{
			DeviceID:     "TEST-DEVICE",
			TradePointID: 2,
		})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		kaspiErr, ok := domain.IsKaspiError(err)
		if !ok {
			t.Fatalf("Expected KaspiError, got %T: %v", err, err)
		}

		if kaspiErr.StatusCode != -1503 {
			t.Errorf("Expected error code -1503, got %d", kaspiErr.StatusCode)
		}
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Run("successfully deletes device", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/device/delete" {
				t.Errorf("Expected URL path /device/delete, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			// Verify request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close()

			var reqBody struct {
				DeviceToken string `json:"DeviceToken"`
			}
			err := json.Unmarshal(body, &reqBody)
			if err != nil {
				t.Errorf("Failed to parse request body: %v", err)
			}

			if reqBody.DeviceToken != "test-token" {
				t.Errorf("Expected DeviceToken test-token, got %s", reqBody.DeviceToken)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK"
			}`), nil
		}

		err := svc.DeleteDevice(context.Background(), "test-token")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles error response", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": -1501,
				"Message": "Device not found"
			}`), nil
		}

		err := svc.DeleteDevice(context.Background(), "invalid-token")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		kaspiErr, ok := domain.IsKaspiError(err)
		if !ok {
			t.Fatalf("Expected KaspiError, got %T", err)
		}

		if kaspiErr.StatusCode != -1501 {
			t.Errorf("Expected status code -1501, got %d", kaspiErr.StatusCode)
		}
	})
}

//////// 	End of device operations testing		////////

//////// 	Payment operations testing		////////

func TestCreateQR(t *testing.T) {
	t.Run("successfully creates QR token", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		expireDate, _ := time.Parse(time.RFC3339, "2023-05-16T10:30:00+06:00")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/qr/create" {
				t.Errorf("Expected URL path /qr/create, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			// Verify request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close()

			var reqBody domain.QRCreateRequest
			err := json.Unmarshal(body, &reqBody)
			if err != nil {
				t.Errorf("Failed to parse request body: %v", err)
			}

			if reqBody.DeviceToken != "test-token" {
				t.Errorf("Expected DeviceToken test-token, got %s", reqBody.DeviceToken)
			}

			if reqBody.Amount != 200.00 {
				t.Errorf("Expected Amount 200.00, got %f", reqBody.Amount)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"QrToken": "51236903777280167836178166503744993984459",
					"ExpireDate": "2023-05-16T10:30:00+06:00",
					"QrPaymentId": 15,
					"PaymentMethods": ["Gold", "Red", "Loan"],
					"QrPaymentBehaviorOptions": {
						"StatusPollingInterval": 5,
						"QrCodeScanWaitTimeout": 180,
						"PaymentConfirmationTimeout": 65
					}
				}
			}`), nil
		}

		resp, err := svc.CreateQR(context.Background(), domain.QRCreateRequest{
			DeviceToken: "test-token",
			Amount:      200.00,
			ExternalID:  "15",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.QrToken != "51236903777280167836178166503744993984459" {
			t.Errorf("Expected QR token 51236903777280167836178166503744993984459, got %s", resp.QrToken)
		}

		if resp.QrPaymentID != 15 {
			t.Errorf("Expected QrPaymentID 15, got %d", resp.QrPaymentID)
		}

		if !resp.ExpireDate.Equal(expireDate) {
			t.Errorf("Expected ExpireDate %v, got %v", expireDate, resp.ExpireDate)
		}

		if len(resp.PaymentMethods) != 3 {
			t.Errorf("Expected 3 payment methods, got %d", len(resp.PaymentMethods))
		}

		if resp.QrPaymentBehaviorOptions.StatusPollingInterval != 5 {
			t.Errorf("Expected StatusPollingInterval 5, got %d", resp.QrPaymentBehaviorOptions.StatusPollingInterval)
		}
	})
}

func TestCreatePaymentLink(t *testing.T) {
	t.Run("successfully creates payment link", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/qr/create-link" {
				t.Errorf("Expected URL path /qr/create-link, got %s", req.URL.Path)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"PaymentLink": "https://pay.kaspi.kz/pay/123456789",
					"ExpireDate": "2023-05-16T10:30:00+06:00",
					"PaymentId": 15,
					"PaymentMethods": ["Gold", "Red", "Loan"],
					"PaymentBehaviorOptions": {
						"StatusPollingInterval": 5,
						"LinkActivationWaitTimeout": 180,
						"PaymentConfirmationTimeout": 65
					}
				}
			}`), nil
		}

		resp, err := svc.CreatePaymentLink(context.Background(), domain.PaymentLinkCreateRequest{
			DeviceToken: "test-token",
			Amount:      200.00,
			ExternalID:  "15",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.PaymentLink != "https://pay.kaspi.kz/pay/123456789" {
			t.Errorf("Expected payment link https://pay.kaspi.kz/pay/123456789, got %s", resp.PaymentLink)
		}

		if resp.PaymentID != 15 {
			t.Errorf("Expected PaymentID 15, got %d", resp.PaymentID)
		}
	})
}

func TestGetPaymentStatus(t *testing.T) {
	t.Run("successfully gets payment status", func(t *testing.T) {
		log := setupTestLogger()
		svc, mockClient := setupTestService(log, "basic")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/payment/status/15" {
				t.Errorf("Expected URL path /payment/status/15, got %s", req.URL.Path)
			}

			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"Status": "Wait",
					"TransactionId": "35134863",
					"LoanOfferName": "Рассрочка 0-0-12",
					"LoanTerm": 12,
					"IsOffer": true,
					"ProductType": "Loan",
					"Amount": 200.00,
					"StoreName": "Store 1",
					"Address": "Test Address",
					"City": "Almaty"
				}
			}`), nil
		}

		resp, err := svc.GetPaymentStatus(context.Background(), 15)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Status != "Wait" {
			t.Errorf("Expected Status Wait, got %s", resp.Status)
		}

		if resp.TransactionID != "35134863" {
			t.Errorf("Expected TransactionID 35134863, got %s", resp.TransactionID)
		}

		if resp.ProductType != "Loan" {
			t.Errorf("Expected ProductType Loan, got %s", resp.ProductType)
		}

		if resp.Amount != 200.00 {
			t.Errorf("Expected Amount 200.00, got %f", resp.Amount)
		}
	})
}

//////// 	End of payment operations testing		////////
