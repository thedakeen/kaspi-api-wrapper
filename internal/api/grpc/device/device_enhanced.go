package device

import (
	"context"
	grpchandler "kaspi-api-wrapper/internal/api/grpc"
	"kaspi-api-wrapper/internal/domain"
	devicev1 "kaspi-api-wrapper/pkg/protos/gen/go/device"
)

// GetTradePointsEnhanced implements kaspiv1.DeviceServiceServer
func (s *serverAPI) GetTradePointsEnhanced(ctx context.Context, req *devicev1.GetTradePointsEnhancedRequest) (*devicev1.GetTradePointsResponse, error) {
	tradePoints, err := s.deviceEnhancedProvider.GetTradePointsEnhanced(ctx, req.OrganizationBin)
	if err != nil {
		s.log.Error("GetTradePointsEnhanced failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
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

// RegisterDeviceEnhanced implements kaspiv1.DeviceServiceServer
func (s *serverAPI) RegisterDeviceEnhanced(ctx context.Context, req *devicev1.RegisterDeviceEnhancedRequest) (*devicev1.RegisterDeviceResponse, error) {
	domainReq := domain.EnhancedDeviceRegisterRequest{
		DeviceID:        req.DeviceId,
		TradePointID:    req.TradepointId,
		OrganizationBin: req.OrganizationBin,
	}

	result, err := s.deviceEnhancedProvider.RegisterDeviceEnhanced(ctx, domainReq)
	if err != nil {
		s.log.Error("RegisterDeviceEnhanced failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &devicev1.RegisterDeviceResponse{
		DeviceToken: result.DeviceToken,
	}, nil
}

// DeleteDeviceEnhanced implements kaspiv1.DeviceServiceServer
func (s *serverAPI) DeleteDeviceEnhanced(ctx context.Context, req *devicev1.DeleteDeviceEnhancedRequest) (*devicev1.DeleteDeviceResponse, error) {
	domainReq := domain.EnhancedDeviceDeleteRequest{
		DeviceToken:     req.DeviceToken,
		OrganizationBin: req.OrganizationBin,
	}

	err := s.deviceEnhancedProvider.DeleteDeviceEnhanced(ctx, domainReq)
	if err != nil {
		s.log.Error("DeleteDeviceEnhanced failed", "error", err.Error())
		return nil, grpchandler.HandleError(err, s.log)
	}

	return &devicev1.DeleteDeviceResponse{}, nil
}
