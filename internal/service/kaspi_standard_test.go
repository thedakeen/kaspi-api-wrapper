package service_test

import (
	"context"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/service"
	"kaspi-api-wrapper/internal/testutils"
	"net/http"
	"testing"
	"time"
)

func TestCreateRefundQR(t *testing.T) {
	t.Run("successfully creates refund QR token", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		expireDate, _ := time.Parse(time.RFC3339, "2023-05-16T10:30:00+06:00")

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/return/create" {
				t.Errorf("Expected URL path /return/create, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"QrToken": "51236903777280167836178166503744993984459",
					"ExpireDate": "2023-05-16T10:30:00+06:00",
					"QrReturnId": 15,
					"QrReturnBehaviorOptions": {
						"QrCodeScanEventPollingInterval": 5,
						"QrCodeScanWaitTimeout": 180
					}
				}
			}`), nil
		}

		resp, err := svc.CreateRefundQR(context.Background(), domain.QRRefundCreateRequest{
			DeviceToken: "test-token",
			ExternalID:  "15",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.QrToken != "51236903777280167836178166503744993984459" {
			t.Errorf("Expected QR token 51236903777280167836178166503744993984459, got %s", resp.QrToken)
		}

		if resp.QrReturnID != 15 {
			t.Errorf("Expected QrReturnID 15, got %d", resp.QrReturnID)
		}

		if !resp.ExpireDate.Equal(expireDate) {
			t.Errorf("Expected ExpireDate %v, got %v", expireDate, resp.ExpireDate)
		}

		if resp.QrRefundBehaviorOptions.QrCodeScanEventPollingInterval != 5 {
			t.Errorf("Expected QrCodeScanEventPollingInterval 5, got %d",
				resp.QrRefundBehaviorOptions.QrCodeScanEventPollingInterval)
		}
	})

	t.Run("fails on basic scheme", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "basic", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		_, err := svc.CreateRefundQR(context.Background(), domain.QRRefundCreateRequest{
			DeviceToken: "test-token",
			ExternalID:  "15",
		})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "refund functionality is not available in basic scheme" {
			t.Errorf("Unexpected error message: %s", err.Error())
		}
	})
}

func TestGetRefundStatus(t *testing.T) {
	t.Run("successfully gets refund status", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/return/status/15" {
				t.Errorf("Expected URL path /return/status/15, got %s", req.URL.Path)
			}

			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"Status": "QrTokenCreated"
				}
			}`), nil
		}

		resp, err := svc.GetRefundStatus(context.Background(), 15)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Status != "QrTokenCreated" {
			t.Errorf("Expected Status QrTokenCreated, got %s", resp.Status)
		}
	})
}

func TestGetCustomerOperations(t *testing.T) {
	t.Run("successfully gets customer operations", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/return/operations" {
				t.Errorf("Expected URL path /return/operations, got %s", req.URL.Path)
			}

			if req.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": [
					{
						"QrPaymentId": 900077110,
						"TransactionDate": "2023-05-16T10:30:00+06:00",
						"Amount": 1.00
					},
					{
						"QrPaymentId": 900077111,
						"TransactionDate": "2023-05-16T10:30:00+06:00",
						"Amount": 2.00
					}
				]
			}`), nil
		}

		operations, err := svc.GetCustomerOperations(context.Background(), domain.CustomerOperationsRequest{
			DeviceToken: "test-token",
			QrReturnID:  15,
			MaxResult:   10,
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(operations))
		}

		if operations[0].QrPaymentID != 900077110 {
			t.Errorf("Expected QrPaymentID 900077110, got %d", operations[0].QrPaymentID)
		}

		if operations[0].Amount != 1.00 {
			t.Errorf("Expected Amount 1.00, got %f", operations[0].Amount)
		}

		if operations[1].QrPaymentID != 900077111 {
			t.Errorf("Expected QrPaymentID 900077111, got %d", operations[1].QrPaymentID)
		}
	})
}

func TestGetPaymentDetails(t *testing.T) {
	t.Run("successfully gets payment details", func(t *testing.T) {
		log := setupTestLogger()
		svc := service.NewKaspiService(
			log, "standard", "https://test.com",
			"https://test.com", "https://test.com", "test-api-key",
		)

		mockClient := &testutils.MockHTTPClient{}
		svc.SetHTTPClient(mockClient)

		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			expectedPath := "/payment/details"
			if req.URL.Path != expectedPath {
				t.Errorf("Expected URL path %s, got %s", expectedPath, req.URL.Path)
			}

			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			if req.URL.Query().Get("QrPaymentId") != "123" {
				t.Errorf("Expected QrPaymentId 123, got %s", req.URL.Query().Get("QrPaymentId"))
			}

			if req.URL.Query().Get("DeviceToken") != "test-token" {
				t.Errorf("Expected DeviceToken test-token, got %s", req.URL.Query().Get("DeviceToken"))
			}

			return testutils.NewMockResponse(http.StatusOK, `{
				"StatusCode": 0,
				"Message": "OK",
				"Data": {
					"QrPaymentId": 123,
					"TotalAmount": 11.00,
					"AvailableReturnAmount": 11.00,
					"TransactionDate": "2021-11-03T11:55:14.166+06:00"
				}
			}`), nil
		}

		details, err := svc.GetPaymentDetails(context.Background(), 123, "test-token")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
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
}

func TestRefundPayment(t *testing.T) {
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

		resp, err := svc.RefundPayment(context.Background(), domain.RefundRequest{
			DeviceToken: "test-token",
			QrPaymentID: 123,
			QrReturnID:  13,
			Amount:      10.00,
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.ReturnOperationID != 15 {
			t.Errorf("Expected ReturnOperationID 15, got %d", resp.ReturnOperationID)
		}
	})
}
