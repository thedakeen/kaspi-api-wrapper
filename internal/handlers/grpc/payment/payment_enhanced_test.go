package payment_test

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	paymentv1 "kaspi-api-wrapper/pkg/protos/gen/go/payment"
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

func TestCreateQREnhanced(t *testing.T) {
	t.Run("successfully creates enhanced QR", func(t *testing.T) {
		expireDate := time.Now().Add(5 * time.Minute)

		mockEnhancedProvider := &MockPaymentEnhancedProvider{
			CreateQREnhancedFunc: func(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}
				if req.Amount != 200.00 {
					return nil, &domain.KaspiError{StatusCode: -990000028, Message: "Invalid payment amount"}
				}
				if req.OrganizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -99000002, Message: "Organization not found"}
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

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &paymentv1.CreateQREnhancedRequest{
			DeviceToken:     "test-token",
			Amount:          200.00,
			ExternalId:      "15",
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.CreateQREnhanced(context.Background(), req)

		if err != nil {
			t.Fatalf("CreateQREnhanced returned error: %v", err)
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

	t.Run("handles invalid payment amount error", func(t *testing.T) {
		mockEnhancedProvider := &MockPaymentEnhancedProvider{
			CreateQREnhancedFunc: func(ctx context.Context, req domain.EnhancedQRCreateRequest) (*domain.QRCreateResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: 990000028,
					Message:    "Invalid payment amount",
				}
			},
		}

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &paymentv1.CreateQREnhancedRequest{
			DeviceToken:     "test-token",
			Amount:          -200.00, // Invalid amount
			ExternalId:      "15",
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.CreateQREnhanced(context.Background(), req)

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

		if st.Message() != "Invalid payment amount" {
			t.Errorf("Expected message 'Invalid payment amount', got %q", st.Message())
		}
	})
}

func TestCreatePaymentLinkEnhanced(t *testing.T) {
	t.Run("successfully creates enhanced payment link", func(t *testing.T) {
		expireDate := time.Now().Add(5 * time.Minute)

		mockEnhancedProvider := &MockPaymentEnhancedProvider{
			CreatePaymentLinkEnhancedFunc: func(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}
				if req.Amount != 200.00 {
					return nil, &domain.KaspiError{StatusCode: -990000028, Message: "Invalid payment amount"}
				}
				if req.OrganizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -99000002, Message: "Organization not found"}
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

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &paymentv1.CreatePaymentLinkEnhancedRequest{
			DeviceToken:     "test-token",
			Amount:          200.00,
			ExternalId:      "15",
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.CreatePaymentLinkEnhanced(context.Background(), req)

		if err != nil {
			t.Fatalf("CreatePaymentLinkEnhanced returned error: %v", err)
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

	t.Run("handles trade point disabled error", func(t *testing.T) {
		mockEnhancedProvider := &MockPaymentEnhancedProvider{
			CreatePaymentLinkEnhancedFunc: func(ctx context.Context, req domain.EnhancedPaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: 990000018,
					Message:    "Trade point is disabled",
				}
			},
		}

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &paymentv1.CreatePaymentLinkEnhancedRequest{
			DeviceToken:     "test-token",
			Amount:          200.00,
			ExternalId:      "15",
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.CreatePaymentLinkEnhanced(context.Background(), req)

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

		if st.Code() != codes.FailedPrecondition {
			t.Errorf("Expected status code %s, got %s", codes.FailedPrecondition, st.Code())
		}

		if st.Message() != "Trade point is disabled" {
			t.Errorf("Expected message 'Trade point is disabled', got %q", st.Message())
		}
	})
}
