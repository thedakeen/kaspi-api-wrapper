package refund_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/handlers/grpc/refund"
	refundv1 "kaspi-api-wrapper/pkg/protos/gen/go/refund"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

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

func createTestServer(refundProvider *MockRefundProvider) *refundServer {
	log := setupTestLogger()
	srv := &refundServer{
		server: refund.RegisterTest(log, refundProvider),
	}
	return srv
}

type refundServer struct {
	server refundv1.RefundServiceServer
}

func TestCreateRefundQR(t *testing.T) {
	t.Run("successfully creates refund QR", func(t *testing.T) {
		expireDate := time.Now().Add(5 * time.Minute)

		mockProvider := &MockRefundProvider{
			CreateRefundQRFunc: func(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
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

		srv := createTestServer(mockProvider)
		req := &refundv1.CreateRefundQRRequest{
			DeviceToken: "test-token",
			ExternalId:  "15",
		}

		resp, err := srv.server.CreateRefundQR(context.Background(), req)

		if err != nil {
			t.Fatalf("CreateRefundQR returned error: %v", err)
		}

		if resp.QrToken != "51236903777280167836178166503744993984459" {
			t.Errorf("Expected QR token 51236903777280167836178166503744993984459, got %s", resp.QrToken)
		}

		if resp.QrReturnId != 15 {
			t.Errorf("Expected QR return ID 15, got %d", resp.QrReturnId)
		}

		if resp.QrRefundBehaviorOptions.QrCodeScanEventPollingInterval != 5 {
			t.Errorf("Expected QrCodeScanEventPollingInterval 5, got %d",
				resp.QrRefundBehaviorOptions.QrCodeScanEventPollingInterval)
		}
	})

	t.Run("handles invalid request", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			CreateRefundQRFunc: func(ctx context.Context, req domain.QRRefundCreateRequest) (*domain.QRRefundCreateResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -1501,
					Message:    "Device not found",
				}
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundv1.CreateRefundQRRequest{
			DeviceToken: "invalid-token",
		}

		resp, err := srv.server.CreateRefundQR(context.Background(), req)

		if resp != nil {
			t.Errorf("Expected nil response, got %+v", resp)
		}

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("Expected gRPC status error, got %T: %v", err, err)
		}

		if st.Code() != codes.NotFound {
			t.Errorf("Expected status code %s, got %s", codes.NotFound, st.Code())
		}
	})
}

func TestGetRefundStatus(t *testing.T) {
	t.Run("successfully gets refund status", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			GetRefundStatusFunc: func(ctx context.Context, qrReturnID int64) (*domain.RefundStatusResponse, error) {
				if qrReturnID != 15 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Refund not found"}
				}

				return &domain.RefundStatusResponse{
					Status: "QrTokenCreated",
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundv1.GetRefundStatusRequest{
			QrReturnId: 15,
		}

		resp, err := srv.server.GetRefundStatus(context.Background(), req)

		if err != nil {
			t.Fatalf("GetRefundStatus returned error: %v", err)
		}

		if resp.Status != "QrTokenCreated" {
			t.Errorf("Expected status QrTokenCreated, got %s", resp.Status)
		}
	})
}

func TestGetCustomerOperations(t *testing.T) {
	t.Run("successfully gets customer operations", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			GetCustomerOperationsFunc: func(ctx context.Context, req domain.CustomerOperationsRequest) ([]domain.CustomerOperation, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				if req.QrReturnID != 15 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Refund not found"}
				}

				transactionDate := time.Now().Add(-24 * time.Hour)

				return []domain.CustomerOperation{
					{
						QrPaymentID:     900077110,
						TransactionDate: transactionDate,
						Amount:          100.00,
					},
					{
						QrPaymentID:     900077111,
						TransactionDate: transactionDate,
						Amount:          200.00,
					},
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundv1.GetCustomerOperationsRequest{
			DeviceToken: "test-token",
			QrReturnId:  15,
			MaxResult:   10,
		}

		resp, err := srv.server.GetCustomerOperations(context.Background(), req)

		if err != nil {
			t.Fatalf("GetCustomerOperations returned error: %v", err)
		}

		if len(resp.Operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(resp.Operations))
		}

		if resp.Operations[0].QrPaymentId != 900077110 {
			t.Errorf("Expected QrPaymentId 900077110, got %d", resp.Operations[0].QrPaymentId)
		}

		if resp.Operations[0].Amount != 100.00 {
			t.Errorf("Expected amount 100.00, got %f", resp.Operations[0].Amount)
		}

		if resp.Operations[1].QrPaymentId != 900077111 {
			t.Errorf("Expected QrPaymentId 900077111, got %d", resp.Operations[1].QrPaymentId)
		}
	})
}

func TestGetPaymentDetails(t *testing.T) {
	t.Run("successfully gets payment details", func(t *testing.T) {
		transactionDate := time.Now().Add(-24 * time.Hour)

		mockProvider := &MockRefundProvider{
			GetPaymentDetailsFunc: func(ctx context.Context, qrPaymentID int64, deviceToken string) (*domain.PaymentDetailsResponse, error) {
				if qrPaymentID != 123 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"}
				}

				if deviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				return &domain.PaymentDetailsResponse{
					QrPaymentID:           123,
					TotalAmount:           100.00,
					AvailableReturnAmount: 100.00,
					TransactionDate:       transactionDate,
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundv1.GetPaymentDetailsRequest{
			QrPaymentId: 123,
			DeviceToken: "test-token",
		}

		resp, err := srv.server.GetPaymentDetails(context.Background(), req)

		if err != nil {
			t.Fatalf("GetPaymentDetails returned error: %v", err)
		}

		if resp.QrPaymentId != 123 {
			t.Errorf("Expected QrPaymentId 123, got %d", resp.QrPaymentId)
		}

		if resp.TotalAmount != 100.00 {
			t.Errorf("Expected total amount 100.00, got %f", resp.TotalAmount)
		}

		if resp.AvailableReturnAmount != 100.00 {
			t.Errorf("Expected available return amount 100.00, got %f", resp.AvailableReturnAmount)
		}
	})
}

func TestRefundPayment(t *testing.T) {
	t.Run("successfully refunds payment", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			RefundPaymentFunc: func(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				if req.QrPaymentID != 123 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"}
				}

				if req.QrReturnID != 15 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Refund not found"}
				}

				if req.Amount <= 0 {
					return nil, &domain.KaspiError{StatusCode: -99000005, Message: "Invalid amount"}
				}

				return &domain.RefundResponse{
					ReturnOperationID: 20,
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundv1.RefundPaymentRequest{
			DeviceToken: "test-token",
			QrPaymentId: 123,
			QrReturnId:  15,
			Amount:      50.00,
		}

		resp, err := srv.server.RefundPayment(context.Background(), req)

		if err != nil {
			t.Fatalf("RefundPayment returned error: %v", err)
		}

		if resp.ReturnOperationId != 20 {
			t.Errorf("Expected ReturnOperationId 20, got %d", resp.ReturnOperationId)
		}
	})

	t.Run("handles refund amount exceeds purchase amount error", func(t *testing.T) {
		mockProvider := &MockRefundProvider{
			RefundPaymentFunc: func(ctx context.Context, req domain.RefundRequest) (*domain.RefundResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -99000005,
					Message:    "Refund amount cannot exceed the purchase amount",
				}
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundv1.RefundPaymentRequest{
			DeviceToken: "test-token",
			QrPaymentId: 123,
			QrReturnId:  15,
			Amount:      1000.00,
		}

		resp, err := srv.server.RefundPayment(context.Background(), req)

		if resp != nil {
			t.Errorf("Expected nil response, got %+v", resp)
		}

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("Expected gRPC status error, got %T: %v", err, err)
		}

		if st.Code() != codes.InvalidArgument {
			t.Errorf("Expected status code %s, got %s", codes.InvalidArgument, st.Code())
		}

		expectedMsg := "Refund amount cannot exceed the purchase amount"
		if st.Message() != expectedMsg {
			t.Errorf("Expected message %q, got %q", expectedMsg, st.Message())
		}
	})
}
