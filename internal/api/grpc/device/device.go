package device

import (
	"context"
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/validator"
	devicev1 "kaspi-api-wrapper/pkg/protos/gen/go/device"
	"log/slog"
)

type serverAPI struct {
	devicev1.UnimplementedDeviceServiceServer
	log                    *slog.Logger
	deviceProvider         api.DeviceProvider
	deviceEnhancedProvider api.DeviceEnhancedProvider
}

func Register(gRPC *grpc.Server, log *slog.Logger, deviceProvider api.DeviceProvider, deviceEnhancedProvider api.DeviceEnhancedProvider) {
	devicev1.RegisterDeviceServiceServer(gRPC, &serverAPI{
		log:                    log,
		deviceProvider:         deviceProvider,
		deviceEnhancedProvider: deviceEnhancedProvider,
	})
}

// GetTradePoints implements kaspiv1.DeviceServiceServer
func (s *serverAPI) GetTradePoints(ctx context.Context, req *devicev1.GetTradePointsRequest) (*devicev1.GetTradePointsResponse, error) {
	log := s.log.With(
		slog.String("method", "GetTradePoints"),
	)

	tradePoints, err := s.deviceProvider.GetTradePoints(ctx)
	if err != nil {
		log.Error("GetTradePoints failed", "error", err.Error())
		return nil, grpchandler.HandleKaspiError(err, log)
	}

	resp := &devicev1.GetTradePointsResponse{
		Tradepoints: make([]*devicev1.TradePoint, 0, len(tradePoints)),
	}

	for _, tp := range tradePoints {
		resp.Tradepoints = append(resp.Tradepoints, &devicev1.TradePoint{
			TradepointId:   tp.TradePointID,
			TradepointName: tp.TradePointName,
		})
	}

	return resp, nil
}

// RegisterDevice implements kaspiv1.DeviceServiceServer
func (s *serverAPI) RegisterDevice(ctx context.Context, req *devicev1.RegisterDeviceRequest) (*devicev1.RegisterDeviceResponse, error) {
	log := s.log.With(
		slog.String("method", "RegisterDevice"),
		slog.String("deviceId", req.DeviceId),
		slog.Int64("tradepointId", req.TradepointId),
	)

	domainReq := domain.DeviceRegisterRequest{
		DeviceID:     req.DeviceId,
		TradePointID: req.TradepointId,
	}

	if err := validator.ValidateDeviceRegisterRequest(domainReq); err != nil {
		return nil, validator.GRPCError(err)
	}

	result, err := s.deviceProvider.RegisterDevice(ctx, domainReq)
	if err != nil {
		log.Error("failed to register device", "error", err.Error())
		return nil, grpchandler.HandleKaspiError(err, log)
	}

	return &devicev1.RegisterDeviceResponse{
		DeviceToken: result.DeviceToken,
	}, nil
}

// DeleteDevice implements kaspiv1.DeviceServiceServer
func (s *serverAPI) DeleteDevice(ctx context.Context, req *devicev1.DeleteDeviceRequest) (*devicev1.DeleteDeviceResponse, error) {
	log := s.log.With(
		slog.String("method", "DeleteDevice"),
		slog.String("deviceToken", req.DeviceToken),
	)

	if err := validator.ValidateDeviceToken(req.DeviceToken); err != nil {
		return nil, validator.GRPCError(err)
	}

	err := s.deviceProvider.DeleteDevice(ctx, req.DeviceToken)
	if err != nil {
		log.Error("failed to delete device", "error", err.Error())
		return nil, grpchandler.HandleKaspiError(err, log)
	}

	return &devicev1.DeleteDeviceResponse{}, nil
}
