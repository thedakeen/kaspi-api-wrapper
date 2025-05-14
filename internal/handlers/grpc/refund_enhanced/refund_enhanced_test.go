package refund_enhanced_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/handlers/grpc/refund_enhanced"
	refundenhancedv1 "kaspi-api-wrapper/pkg/protos/gen/go/refund_enhanced"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

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

func createTestServer(refundEnhancedProvider *MockRefundEnhancedProvider) *refundEnhancedServer {
	log := setupTestLogger()
	srv := &refundEnhancedServer{
		server: refund_enhanced.RegisterTest(log, refundEnhancedProvider),
	}
	return srv
}

type refundEnhancedServer struct {
	server refundenhancedv1.EnhancedRefundServiceServer
}

func TestRefundPaymentEnhanced(t *testing.T) {
	t.Run("successfully refunds payment enhanced", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			RefundPaymentEnhancedFunc: func(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error) {
				if req.DeviceToken != "test-token" {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				if req.QrPaymentID != 123 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"}
				}

				if req.Amount <= 0 {
					return nil, &domain.KaspiError{StatusCode: -99000005, Message: "Invalid amount"}
				}

				if req.OrganizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -99000002, Message: "Organization not found"}
				}

				return &domain.RefundResponse{
					ReturnOperationID: 20,
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.RefundPaymentEnhancedRequest{
			DeviceToken:     "test-token",
			QrPaymentId:     123,
			Amount:          50.00,
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.RefundPaymentEnhanced(context.Background(), req)

		if err != nil {
			t.Fatalf("RefundPaymentEnhanced returned error: %v", err)
		}

		if resp.ReturnOperationId != 20 {
			t.Errorf("Expected ReturnOperationId 20, got %d", resp.ReturnOperationId)
		}
	})

	t.Run("handles missing organization BIN error", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			RefundPaymentEnhancedFunc: func(ctx context.Context, req domain.EnhancedRefundRequest) (*domain.RefundResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -99000002,
					Message:    "Trade point not found",
				}
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.RefundPaymentEnhancedRequest{
			DeviceToken:     "test-token",
			QrPaymentId:     123,
			Amount:          50.00,
			OrganizationBin: "invalid-bin",
		}

		resp, err := srv.server.RefundPaymentEnhanced(context.Background(), req)

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

func TestGetClientInfo(t *testing.T) {
	t.Run("successfully gets client info", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			GetClientInfoFunc: func(ctx context.Context, phoneNumber string, deviceToken int64) (*domain.ClientInfoResponse, error) {
				if phoneNumber != "87071234567" {
					return nil, &domain.KaspiError{StatusCode: -1404, Message: "Invalid phone number"}
				}

				if deviceToken != 2 {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				return &domain.ClientInfoResponse{
					ClientName: "John Doe",
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.GetClientInfoRequest{
			PhoneNumber: "87071234567",
			DeviceToken: 2,
		}

		resp, err := srv.server.GetClientInfo(context.Background(), req)

		if err != nil {
			t.Fatalf("GetClientInfo returned error: %v", err)
		}

		if resp.ClientName != "John Doe" {
			t.Errorf("Expected client name John Doe, got %s", resp.ClientName)
		}
	})

	t.Run("handles invalid phone number error", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			GetClientInfoFunc: func(ctx context.Context, phoneNumber string, deviceToken int64) (*domain.ClientInfoResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -1404,
					Message:    "Invalid phone number",
				}
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.GetClientInfoRequest{
			PhoneNumber: "invalid-phone",
			DeviceToken: 2,
		}

		resp, err := srv.server.GetClientInfo(context.Background(), req)

		if resp != nil {
			t.Errorf("Expected nil response, got %+v", resp)
		}

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestCreateRemotePayment(t *testing.T) {
	t.Run("successfully creates remote payment", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			CreateRemotePaymentFunc: func(ctx context.Context, req domain.RemotePaymentRequest) (*domain.RemotePaymentResponse, error) {
				if req.PhoneNumber != "87071234567" {
					return nil, &domain.KaspiError{StatusCode: -1404, Message: "Invalid phone number"}
				}

				if req.DeviceToken != 2 {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				if req.Amount <= 0 {
					return nil, &domain.KaspiError{StatusCode: -990000028, Message: "Invalid payment amount"}
				}

				if req.OrganizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -99000002, Message: "Organization not found"}
				}

				return &domain.RemotePaymentResponse{
					QrPaymentID: 15,
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.CreateRemotePaymentRequest{
			OrganizationBin: "180340021791",
			Amount:          100.00,
			PhoneNumber:     "87071234567",
			DeviceToken:     "2",
			Comment:         "Test payment",
		}

		resp, err := srv.server.CreateRemotePayment(context.Background(), req)

		if err != nil {
			t.Fatalf("CreateRemotePayment returned error: %v", err)
		}

		if resp.QrPaymentId != 15 {
			t.Errorf("Expected QrPaymentId 15, got %d", resp.QrPaymentId)
		}
	})

	t.Run("handles invalid device token format", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.CreateRemotePaymentRequest{
			OrganizationBin: "180340021791",
			Amount:          100.00,
			PhoneNumber:     "87071234567",
			DeviceToken:     "invalid",
			Comment:         "Test payment",
		}

		resp, err := srv.server.CreateRemotePayment(context.Background(), req)

		if resp != nil {
			t.Errorf("Expected nil response, got %+v", resp)
		}

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestCancelRemotePayment(t *testing.T) {
	t.Run("successfully cancels remote payment", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			CancelRemotePaymentFunc: func(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error) {
				if req.DeviceToken != 2 {
					return nil, &domain.KaspiError{StatusCode: -1501, Message: "Device not found"}
				}

				if req.QrPaymentID != 15 {
					return nil, &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"}
				}

				if req.OrganizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -99000002, Message: "Organization not found"}
				}

				return &domain.RemotePaymentCancelResponse{
					Status: "RemotePaymentCanceled",
				}, nil
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.CancelRemotePaymentRequest{
			OrganizationBin: "180340021791",
			QrPaymentId:     15,
			DeviceToken:     2,
		}

		resp, err := srv.server.CancelRemotePayment(context.Background(), req)

		if err != nil {
			t.Fatalf("CancelRemotePayment returned error: %v", err)
		}

		if resp.Status != "RemotePaymentCanceled" {
			t.Errorf("Expected status RemotePaymentCanceled, got %s", resp.Status)
		}
	})

	t.Run("handles payment not found error", func(t *testing.T) {
		mockProvider := &MockRefundEnhancedProvider{
			CancelRemotePaymentFunc: func(ctx context.Context, req domain.RemotePaymentCancelRequest) (*domain.RemotePaymentCancelResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -99000001,
					Message:    "Payment with the specified ID not found",
				}
			},
		}

		srv := createTestServer(mockProvider)
		req := &refundenhancedv1.CancelRemotePaymentRequest{
			OrganizationBin: "180340021791",
			QrPaymentId:     999,
			DeviceToken:     2,
		}

		resp, err := srv.server.CancelRemotePayment(context.Background(), req)

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

		expectedMsg := "Payment with the specified ID not found"
		if st.Message() != expectedMsg {
			t.Errorf("Expected message %q, got %q", expectedMsg, st.Message())
		}
	})
}
