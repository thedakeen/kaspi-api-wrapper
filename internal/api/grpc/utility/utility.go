package utility

import (
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	utilityv1 "kaspi-api-wrapper/pkg/protos/gen/go/utility"
)

type serverAPI struct {
	utilityv1.UnimplementedUtilityServiceServer
	utilityProvider api.UtilityProvider
}

func Register(gRPC *grpc.Server, utilityProvider api.UtilityProvider) {
	utilityv1.RegisterUtilityServiceServer(gRPC, &serverAPI{
		utilityProvider: utilityProvider,
	})
}
