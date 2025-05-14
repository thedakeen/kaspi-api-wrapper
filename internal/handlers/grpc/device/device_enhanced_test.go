package device_test

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	devicev1 "kaspi-api-wrapper/pkg/protos/gen/go/device"
	"testing"
)

type MockDeviceEnhancedProvider struct {
	GetTradePointsEnhancedFunc func(ctx context.Context, organizationBin string) ([]domain.TradePoint, error)
	RegisterDeviceEnhancedFunc func(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error)
	DeleteDeviceEnhancedFunc   func(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error
}

func (m *MockDeviceEnhancedProvider) GetTradePointsEnhanced(ctx context.Context, organizationBin string) ([]domain.TradePoint, error) {
	return m.GetTradePointsEnhancedFunc(ctx, organizationBin)
}

func (m *MockDeviceEnhancedProvider) RegisterDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
	return m.RegisterDeviceEnhancedFunc(ctx, req)
}

func (m *MockDeviceEnhancedProvider) DeleteDeviceEnhanced(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error {
	return m.DeleteDeviceEnhancedFunc(ctx, req)
}

func TestGetTradePointsEnhanced(t *testing.T) {
	t.Run("successfully gets enhanced trade points", func(t *testing.T) {
		mockEnhancedProvider := &MockDeviceEnhancedProvider{
			GetTradePointsEnhancedFunc: func(ctx context.Context, organizationBin string) ([]domain.TradePoint, error) {
				if organizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Wrong organization BIN"}
				}

				return []domain.TradePoint{
					{TradePointID: 1, TradePointName: "Store 1"},
					{TradePointID: 2, TradePointName: "Store 2"},
				}, nil
			},
		}

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &devicev1.GetTradePointsEnhancedRequest{
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.GetTradePointsEnhanced(context.Background(), req)

		if err != nil {
			t.Fatalf("GetTradePointsEnhanced returned error: %v", err)
		}

		if len(resp.Tradepoints) != 2 {
			t.Errorf("Expected 2 trade points, got %d", len(resp.Tradepoints))
		}

		if resp.Tradepoints[0].TradepointId != 1 || resp.Tradepoints[0].TradepointName != "Store 1" {
			t.Errorf("Unexpected first trade point: %+v", resp.Tradepoints[0])
		}
	})

	t.Run("handles invalid organization BIN error", func(t *testing.T) {
		mockEnhancedProvider := &MockDeviceEnhancedProvider{
			GetTradePointsEnhancedFunc: func(ctx context.Context, organizationBin string) ([]domain.TradePoint, error) {
				return nil, &domain.KaspiError{
					StatusCode: -99000002,
					Message:    "Trade point not found",
				}
			},
		}

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &devicev1.GetTradePointsEnhancedRequest{
			OrganizationBin: "invalid-bin",
		}

		resp, err := srv.server.GetTradePointsEnhanced(context.Background(), req)

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

func TestRegisterDeviceEnhanced(t *testing.T) {
	t.Run("successfully registers enhanced device", func(t *testing.T) {
		mockEnhancedProvider := &MockDeviceEnhancedProvider{
			RegisterDeviceEnhancedFunc: func(ctx context.Context, req domain.EnhancedDeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
				if req.DeviceID != "TEST-DEVICE" {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Wrong device ID"}
				}
				if req.TradePointID != 1 {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Wrong trade point ID"}
				}
				if req.OrganizationBin != "180340021791" {
					return nil, &domain.KaspiError{StatusCode: -1, Message: "Wrong organization BIN"}
				}

				return &domain.DeviceRegisterResponse{
					DeviceToken: "2be4cc91-5895-48f8-8bc2-86c7bd419b3b",
				}, nil
			},
		}

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &devicev1.RegisterDeviceEnhancedRequest{
			DeviceId:        "TEST-DEVICE",
			TradepointId:    1,
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.RegisterDeviceEnhanced(context.Background(), req)

		if err != nil {
			t.Fatalf("RegisterDeviceEnhanced returned error: %v", err)
		}

		if resp.DeviceToken != "2be4cc91-5895-48f8-8bc2-86c7bd419b3b" {
			t.Errorf("Expected device token 2be4cc91-5895-48f8-8bc2-86c7bd419b3b, got %s", resp.DeviceToken)
		}
	})
}

func TestDeleteDeviceEnhanced(t *testing.T) {
	t.Run("successfully deletes enhanced device", func(t *testing.T) {
		mockEnhancedProvider := &MockDeviceEnhancedProvider{
			DeleteDeviceEnhancedFunc: func(ctx context.Context, req domain.EnhancedDeviceDeleteRequest) error {
				if req.DeviceToken != "test-token" {
					return &domain.KaspiError{StatusCode: -1, Message: "Wrong device token"}
				}
				if req.OrganizationBin != "180340021791" {
					return &domain.KaspiError{StatusCode: -1, Message: "Wrong organization BIN"}
				}
				return nil
			},
		}

		srv := createTestServer(nil, mockEnhancedProvider)
		req := &devicev1.DeleteDeviceEnhancedRequest{
			DeviceToken:     "test-token",
			OrganizationBin: "180340021791",
		}

		resp, err := srv.server.DeleteDeviceEnhanced(context.Background(), req)

		if err != nil {
			t.Fatalf("DeleteDeviceEnhanced returned error: %v", err)
		}

		if resp == nil {
			t.Error("Expected non-nil response")
		}
	})
}
