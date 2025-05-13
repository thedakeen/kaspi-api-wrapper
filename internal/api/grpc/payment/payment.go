package payment

import (
	"google.golang.org/grpc"
	"kaspi-api-wrapper/internal/api"
	paymentv1 "kaspi-api-wrapper/pkg/protos/gen/go/payment"
)

type serverAPI struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentProvider         api.PaymentProvider
	paymentEnhancedProvider api.PaymentEnhancedProvider
}

func Register(gRPC *grpc.Server, paymentProvider api.PaymentProvider, paymentEnhancedProvider api.PaymentEnhancedProvider) {
	paymentv1.RegisterPaymentServiceServer(gRPC, &serverAPI{
		paymentProvider:         paymentProvider,
		paymentEnhancedProvider: paymentEnhancedProvider,
	})
}
