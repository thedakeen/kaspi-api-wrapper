package device_test

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/handlers/grpc/device"
	devicev1 "kaspi-api-wrapper/pkg/protos/gen/go/device"
	"log/slog"
	"os"
	"testing"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

type MockDeviceProvider struct {
	GetTradePointsFunc func(ctx context.Context) ([]domain.TradePoint, error)
	RegisterDeviceFunc func(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDeviceFunc   func(ctx context.Context, deviceToken string) error
}

func (m *MockDeviceProvider) GetTradePoints(ctx context.Context) ([]domain.TradePoint, error) {
	return m.GetTradePointsFunc(ctx)
}

func (m *MockDeviceProvider) RegisterDevice(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
	return m.RegisterDeviceFunc(ctx, req)
}

func (m *MockDeviceProvider) DeleteDevice(ctx context.Context, deviceToken string) error {
	return m.DeleteDeviceFunc(ctx, deviceToken)
}

func createTestServer(deviceProvider *MockDeviceProvider, deviceEnhancedProvider *MockDeviceEnhancedProvider) *deviceServer {
	log := setupTestLogger()
	srv := &deviceServer{
		server: device.RegisterTest(log, deviceProvider, deviceEnhancedProvider),
	}
	return srv
}

type deviceServer struct {
	server devicev1.DeviceServiceServer
}

func TestGetTradePoints(t *testing.T) {
	t.Run("successfully gets trade points", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			GetTradePointsFunc: func(ctx context.Context) ([]domain.TradePoint, error) {
				return []domain.TradePoint{
					{TradePointID: 1, TradePointName: "Store 1"},
					{TradePointID: 2, TradePointName: "Store 2"},
				}, nil
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &devicev1.GetTradePointsRequest{}

		resp, err := srv.server.GetTradePoints(context.Background(), req)

		if err != nil {
			t.Fatalf("GetTradePoints returned error: %v", err)
		}

		if len(resp.Tradepoints) != 2 {
			t.Errorf("Expected 2 trade points, got %d", len(resp.Tradepoints))
		}

		if resp.Tradepoints[0].TradepointId != 1 || resp.Tradepoints[0].TradepointName != "Store 1" {
			t.Errorf("Unexpected first trade point: %+v", resp.Tradepoints[0])
		}

		if resp.Tradepoints[1].TradepointId != 2 || resp.Tradepoints[1].TradepointName != "Store 2" {
			t.Errorf("Unexpected second trade point: %+v", resp.Tradepoints[1])
		}
	})

	t.Run("handles error", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			GetTradePointsFunc: func(ctx context.Context) ([]domain.TradePoint, error) {
				return nil, &domain.KaspiError{
					StatusCode: -14000002,
					Message:    "No trade points available",
				}
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &devicev1.GetTradePointsRequest{}

		resp, err := srv.server.GetTradePoints(context.Background(), req)

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

		expectedMsg := "No trade points available. Please create a trade point in the Kaspi Pay application"
		if st.Message() != expectedMsg {
			t.Errorf("Expected message %q, got %q", expectedMsg, st.Message())
		}
	})
}

func TestRegisterDevice(t *testing.T) {
	t.Run("successfully registers device", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			RegisterDeviceFunc: func(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
				if req.DeviceID != "TEST-DEVICE" {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Wrong device ID"}
				}
				if req.TradePointID != 1 {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Wrong trade point ID"}
				}

				return &domain.DeviceRegisterResponse{
					DeviceToken: "2be4cc91-5895-48f8-8bc2-86c7bd419b3b",
				}, nil
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &devicev1.RegisterDeviceRequest{
			DeviceId:     "TEST-DEVICE",
			TradepointId: 1,
		}

		resp, err := srv.server.RegisterDevice(context.Background(), req)

		if err != nil {
			t.Fatalf("RegisterDevice returned error: %v", err)
		}

		if resp.DeviceToken != "2be4cc91-5895-48f8-8bc2-86c7bd419b3b" {
			t.Errorf("Expected device token 2be4cc91-5895-48f8-8bc2-86c7bd419b3b, got %s", resp.DeviceToken)
		}
	})

	t.Run("handles validation error", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			RegisterDeviceFunc: func(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
				return nil, &domain.KaspiError{
					StatusCode: -1503,
					Message:    "Device is already added to another trade point",
				}
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &devicev1.RegisterDeviceRequest{
			DeviceId:     "TEST-DEVICE",
			TradepointId: 1,
		}

		resp, err := srv.server.RegisterDevice(context.Background(), req)

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

		if st.Code() != codes.AlreadyExists {
			t.Errorf("Expected status code %s, got %s", codes.AlreadyExists, st.Code())
		}
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Run("successfully deletes device", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			DeleteDeviceFunc: func(ctx context.Context, deviceToken string) error {
				if deviceToken != "test-token" {
					return &domain.KaspiError{StatusCode: -1, Message: "Wrong device token"}
				}
				return nil
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &devicev1.DeleteDeviceRequest{
			DeviceToken: "test-token",
		}

		resp, err := srv.server.DeleteDevice(context.Background(), req)

		if err != nil {
			t.Fatalf("DeleteDevice returned error: %v", err)
		}

		if resp == nil {
			t.Error("Expected non-nil response")
		}
	})

	t.Run("handles device not found error", func(t *testing.T) {
		mockProvider := &MockDeviceProvider{
			DeleteDeviceFunc: func(ctx context.Context, deviceToken string) error {
				return &domain.KaspiError{
					StatusCode: -1501,
					Message:    "Device not found",
				}
			},
		}

		srv := createTestServer(mockProvider, nil)
		req := &devicev1.DeleteDeviceRequest{
			DeviceToken: "non-existent-token",
		}

		resp, err := srv.server.DeleteDevice(context.Background(), req)

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

		if st.Message() != "Device not found" {
			t.Errorf("Expected message 'Device not found', got %q", st.Message())
		}
	})
}
