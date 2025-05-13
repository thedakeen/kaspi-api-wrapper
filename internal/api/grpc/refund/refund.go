package refund

import (
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	refundv1 "kaspi-api-wrapper/pkg/protos/gen/go/refund"
)

type serverAPI struct {
	refundv1.UnimplementedRefundServiceServer
	refundProvider api.RefundProvider
}

func Register(gRPC *grpc.Server, refundProvider api.RefundProvider) {
	refundv1.RegisterRefundServiceServer(gRPC, &serverAPI{
		refundProvider: refundProvider,
	})
}
