package payment_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/handlers/grpc/payment"
	paymentv1 "kaspi-api-wrapper/pkg/protos/gen/go/payment"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

type MockPaymentProvider struct {
	CreateQRFunc          func(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error)
	CreatePaymentLinkFunc func(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error)
	GetPaymentStatusFunc  func(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error)
}

func (m *MockPaymentProvider) CreateQR(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error) {
	return m.CreateQRFunc(ctx, req)
}

func (m *MockPaymentProvider) CreatePaymentLink(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
	return m.CreatePaymentLinkFunc(ctx, req)
}

func (m *MockPaymentProvider) GetPaymentStatus(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error) {
	return m.GetPaymentStatusFunc(ctx, qrPaymentID)
}

func createTestServer(paymentProvider *MockPaymentProvider, paymentEnhancedProvider *MockPaymentEnhancedProvider) *paymentServer {
	log := setupTestLogger()
	srv := &paymentServer{
		server: payment.RegisterTest(log, paymentProvider, paymentEnhancedProvider),
	}
	return srv
}

type paymentServer struct {
	server paymentv1.PaymentServiceServer
}

func TestCreateQR(t *testing.T) {
	t.Run("successfully creates QR", func(t *testing.T) {
		expireDate := time.Now().Add(5 * time.Minute)

		mockProvider := &MockPaymentProvider{
			CreateQRFunc: func(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}
				if req.Amount != 200.00 {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Invalid amount"}
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

		srv := createTestServer(mockProvider, nil)
		req := &paymentv1.CreateQRRequest{
			DeviceToken: "test-token",
			Amount:      200.00,
			ExternalId:  "15",
		}

		resp, err := srv.server.CreateQR(context.Background(), req)

		if err != nil {
			t.Fatalf("CreateQR returned error: %v", err)
		}

		if resp.QrToken != "51236903777280167836178166503744993984459" {
			t.Errorf("Expected QR token 51236903777280167836178166503744993984459, got %s", resp.QrToken)
		}

		if resp.QrPaymentId != 15 {
			t.Errorf("Expected QR payment ID 15, got %d", resp.QrPaymentId)
		}

		if len(resp.PaymentMethods) != 3 || resp.PaymentMethods[0] != "Gold" {
			t.Errorf("Unexpected payment methods: %v", resp.PaymentMethods)
		}

		if resp.QrPaymentBehaviorOptions.StatusPollingInterval != 5 {
			t.Errorf("Expected StatusPollingInterval 5, got %d", resp.QrPaymentBehaviorOptions.StatusPollingInterval)
		}
	})

	t.Run("handles device not found error", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{
			CreateQRFunc: func(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -1501,
					Message:    "Device not found",
				}
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &paymentv1.CreateQRRequest{
			DeviceToken: "non-existent-token",
			Amount:      200.00,
		}

		resp, err := srv.server.CreateQR(context.Background(), req)

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

func TestCreatePaymentLink(t *testing.T) {
	t.Run("successfully creates payment link", func(t *testing.T) {
		expireDate := time.Now().Add(5 * time.Minute)

		mockProvider := &MockPaymentProvider{
			CreatePaymentLinkFunc: func(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}
				if req.Amount != 200.00 {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Invalid amount"}
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

		srv := createTestServer(mockProvider, nil)
		req := &paymentv1.CreatePaymentLinkRequest{
			DeviceToken: "test-token",
			Amount:      200.00,
			ExternalId:  "15",
		}

		resp, err := srv.server.CreatePaymentLink(context.Background(), req)

		if err != nil {
			t.Fatalf("CreatePaymentLink returned error: %v", err)
		}

		if resp.PaymentLink != "https://pay.kaspi.kz/pay/123456789" {
			t.Errorf("Expected payment link https://pay.kaspi.kz/pay/123456789, got %s", resp.PaymentLink)
		}

		if resp.PaymentId != 15 {
			t.Errorf("Expected payment ID 15, got %d", resp.PaymentId)
		}

		if len(resp.PaymentMethods) != 3 || resp.PaymentMethods[0] != "Gold" {
			t.Errorf("Unexpected payment methods: %v", resp.PaymentMethods)
		}

		if resp.PaymentBehaviorOptions.StatusPollingInterval != 5 {
			t.Errorf("Expected StatusPollingInterval 5, got %d", resp.PaymentBehaviorOptions.StatusPollingInterval)
		}
	})
}

func TestGetPaymentStatus(t *testing.T) {
	t.Run("successfully gets payment status", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{
			GetPaymentStatusFunc: func(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error) {
				if qrPaymentID != 15 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"}
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

		srv := createTestServer(mockProvider, nil)
		req := &paymentv1.GetPaymentStatusRequest{
			QrPaymentId: 15,
		}

		resp, err := srv.server.GetPaymentStatus(context.Background(), req)

		if err != nil {
			t.Fatalf("GetPaymentStatus returned error: %v", err)
		}

		if resp.Status != "Wait" {
			t.Errorf("Expected status Wait, got %s", resp.Status)
		}

		if resp.TransactionId != "35134863" {
			t.Errorf("Expected transaction ID 35134863, got %s", resp.TransactionId)
		}

		if resp.LoanTerm != 12 {
			t.Errorf("Expected loan term 12, got %d", resp.LoanTerm)
		}

		if resp.ProductType != "Loan" {
			t.Errorf("Expected product type Loan, got %s", resp.ProductType)
		}

		if resp.Amount != 200.00 {
			t.Errorf("Expected amount 200.00, got %f", resp.Amount)
		}
	})

	t.Run("handles payment not found error", func(t *testing.T) {
		mockProvider := &MockPaymentProvider{
			GetPaymentStatusFunc: func(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -1601,
					Message:    "Payment not found",
				}
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &paymentv1.GetPaymentStatusRequest{
			QrPaymentId: 999,
		}

		resp, err := srv.server.GetPaymentStatus(context.Background(), req)

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
