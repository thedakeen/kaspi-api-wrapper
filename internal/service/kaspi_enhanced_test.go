package service_test

import (
	"context"
	"encoding/json"
	"io"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/service"
	"kaspi-api-wrapper/internal/testutils"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGetTradePointsEnhanced(t *testing.T) {
	t.Run("successfully gets trade points", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/partner/tradepoints/180340021791" {
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

		tradePoints, err := svc.GetTradePointsEnhanced(context.Background(), "180340021791")

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
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

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

func TestRegisterDeviceEnhanced(t *testing.T) {
	t.Run("successfully registers device", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/device/register" {
				t.Errorf("Expected URL path /device/register, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			// Verify request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close()

			var reqBody domain.EnhancedDeviceRegisterRequest
			err := json.Unmarshal(body, &reqBody)
			if err != nil {
				t.Errorf("Failed to parse request body: %v", err)
			}

			if reqBody.DeviceID != "TEST-DEVICE" {
				t.Errorf("Expected DeviceID TEST-DEVICE, got %s", reqBody.DeviceID)
			}

			if reqBody.TradePointID != 1 {
				t.Errorf("Expected TradePointID 1, got %d", reqBody.TradePointID)
			}

			if reqBody.OrganizationBin != "180340021791" {
				t.Errorf("Expected OrganizationBin 180340021791, got %s", reqBody.OrganizationBin)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"DeviceToken": "2be4cc91-5895-48f8-8bc2-86c7bd419b3b"
				}
			}`), nil
		}

		resp, err := svc.RegisterDeviceEnhanced(context.Background(), domain.EnhancedDeviceRegisterRequest{
			DeviceID:        "TEST-DEVICE",
			TradePointID:    1,
			OrganizationBin: "180340021791",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.DeviceToken != "2be4cc91-5895-48f8-8bc2-86c7bd419b3b" {
			t.Errorf("Expected device token 2be4cc91-5895-48f8-8bc2-86c7bd419b3b, got %s", resp.DeviceToken)
		}
	})
}

func TestDeleteDeviceEnhanced(t *testing.T) {
	t.Run("successfully deletes device", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

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
				DeviceToken     string `json:"DeviceToken"`
				OrganizationBin string `json:"OrganizationBin"`
			}
			err := json.Unmarshal(body, &reqBody)
			if err != nil {
				t.Errorf("Failed to parse request body: %v", err)
			}

			if reqBody.DeviceToken != "test-token" {
				t.Errorf("Expected DeviceToken test-token, got %s", reqBody.DeviceToken)
			}

			if reqBody.OrganizationBin != "180340021791" {
				t.Errorf("Expected OrganizationBin 180340021791, got %s", reqBody.OrganizationBin)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK"
			}`), nil
		}

		err := svc.DeleteDeviceEnhanced(context.Background(), domain.EnhancedDeviceDeleteRequest{
			DeviceToken:     "test-token",
			OrganizationBin: "180340021791",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("handles error response", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": -1501,
				"Message": "Device not found"
			}`), nil
		}

		err := svc.DeleteDeviceEnhanced(context.Background(), domain.EnhancedDeviceDeleteRequest{
			DeviceToken:     "invalid-token",
			OrganizationBin: "180340021791",
		})

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

func TestCreateQREnhanced(t *testing.T) {
	t.Run("successfully creates QR token", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

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

			var reqBody domain.EnhancedQRCreateRequest
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

			if reqBody.OrganizationBin != "180340021791" {
				t.Errorf("Expected OrganizationBin 180340021791, got %s", reqBody.OrganizationBin)
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

		resp, err := svc.CreateQREnhanced(context.Background(), domain.EnhancedQRCreateRequest{
			OrganizationBin: "180340021791",
			DeviceToken:     "test-token",
			Amount:          200.00,
			ExternalID:      "15",
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

func TestRefundPaymentEnhanced(t *testing.T) {
	t.Run("successfully refunds payment", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/payment/return" {
				t.Errorf("Expected URL path /payment/return, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"ReturnOperationId": 15
				}
			}`), nil
		}

		resp, err := svc.RefundPaymentEnhanced(context.Background(), domain.EnhancedRefundRequest{
			DeviceToken:     "test-token",
			OrganizationBin: "180340021791",
			QrPaymentID:     123,
			Amount:          10.00,
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.ReturnOperationID != 15 {
			t.Errorf("Expected ReturnOperationID 15, got %d", resp.ReturnOperationID)
		}
	})
}

func TestGetClientInfo(t *testing.T) {
	t.Run("successfully gets client info", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "enhanced", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			expectedPath := "/remote/client-info"
			if req.URL.Path != expectedPath {
				t.Errorf("Expected URL path %s, got %s", expectedPath, req.URL.Path)
			}

			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			if req.URL.Query().Get("phoneNumber") != "87071234567" {
				t.Errorf("Expected phoneNumber 87071234567, got %s", req.URL.Query().Get("phoneNumber"))
			}

			if req.URL.Query().Get("deviceToken") != "test-token" {
				t.Errorf("Expected deviceToken test-token, got %s", req.URL.Query().Get("deviceToken"))
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"ClientName": "Walter White"
				}
			}`), nil
		}

		info, err := svc.GetClientInfo(context.Background(), "87071234567", "test-token")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if info.ClientName != "Walter White" {
			t.Errorf("Expected ClientName Walter White, got %s", info.ClientName)
		}
	})
}

func TestCreateRemotePayment(t *testing.T) {
	t.Run("successfully creates remote payment", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "enhanced", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/remote/create" {
				t.Errorf("Expected URL path /remote/create, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"QrPaymentId": 15
				}
			}`), nil
		}

		resp, err := svc.CreateRemotePayment(context.Background(), domain.RemotePaymentRequest{
			OrganizationBin: "180340021791",
			Amount:          100.00,
			PhoneNumber:     "87071234567",
			DeviceToken:     "test-token",
			Comment:         "Test payment",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.QrPaymentID != 15 {
			t.Errorf("Expected QrPaymentID 15, got %d", resp.QrPaymentID)
		}
	})

	t.Run("fails on non-enhanced scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		_, err := svc.CreateRemotePayment(context.Background(), domain.RemotePaymentRequest{
			OrganizationBin: "180340021791",
			Amount:          100.00,
			PhoneNumber:     "87071234567",
			DeviceToken:     "test-token",
			Comment:         "Test payment",
		})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "enhanced scheme") {
			t.Errorf("Expected error message to mention enhanced scheme, got: %s", err.Error())
		}
	})
}

func TestCancelRemotePayment(t *testing.T) {
	t.Run("successfully cancels remote payment", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "enhanced", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/remote/cancel" {
				t.Errorf("Expected URL path /remote/cancel, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"Status": "RemotePaymentCanceled"
				}
			}`), nil
		}

		resp, err := svc.CancelRemotePayment(context.Background(), domain.RemotePaymentCancelRequest{
			OrganizationBin: "180340021791",
			QrPaymentID:     15,
			DeviceToken:     "test-token",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Status != "RemotePaymentCanceled" {
			t.Errorf("Expected Status RemotePaymentCanceled, got %s", resp.Status)
		}
	})

	t.Run("fails on non-enhanced scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		_, err := svc.CancelRemotePayment(context.Background(), domain.RemotePaymentCancelRequest{
			OrganizationBin: "180340021791",
			QrPaymentID:     15,
			DeviceToken:     "test-token",
		})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "enhanced scheme") {
			t.Errorf("Expected error message to mention enhanced scheme, got: %s", err.Error())
		}
	})
}
