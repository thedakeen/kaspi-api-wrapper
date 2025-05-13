package refund_enhanced

import (
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	refundenhancedv1 "kaspi-api-wrapper/pkg/protos/gen/go/refund_enhanced"
)

type serverAPI struct {
	refundenhancedv1.UnimplementedEnhancedRefundServiceServer
	refundEnhancedProvider api.RefundEnhancedProvider
}

func Register(gRPC *grpc.Server, refundEnhancedProvider api.RefundEnhancedProvider) {
	refundenhancedv1.RegisterEnhancedRefundServiceServer(gRPC, &serverAPI{
		refundEnhancedProvider: refundEnhancedProvider,
	})
}
