package device

import (
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	devicev1 "kaspi-api-wrapper/pkg/protos/gen/go/device"
)

type serverAPI struct {
	devicev1.UnimplementedDeviceServiceServer
	deviceProvider         api.DeviceProvider
	deviceEnhancedProvider api.DeviceEnhancedProvider
}

func Register(gRPC *grpc.Server, deviceProvider api.DeviceProvider, deviceEnhancedProvider api.DeviceEnhancedProvider) {
	devicev1.RegisterDeviceServiceServer(gRPC, &serverAPI{
		deviceProvider:         deviceProvider,
		deviceEnhancedProvider: deviceEnhancedProvider,
	})
}
