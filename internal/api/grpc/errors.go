package grpchandler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	"log/slog"
)

// HandleKaspiError handles Kaspi API errors and maps them to appropriate gRPC responses
func HandleKaspiError(err error, log *slog.Logger) error {
	if err != nil && errors.Is(err, domain.ErrUnsupportedFeature) {
		log.Error("scheme compatibility error", "error", err)
		return status.Error(codes.PermissionDenied, err.Error())
	}

	kaspiErr, ok := domain.IsKaspiError(err)
	if !ok {
		log.Error("unexpected error", "error", err)
		return status.Error(codes.Internal, "Internal server error")

	}

	log.Error("kaspi API error",
		"status_code", kaspiErr.StatusCode,
		"message", kaspiErr.Message)

	switch kaspiErr.StatusCode {
	case -1501:
		// Device with the specified identifier not found
		return status.Error(codes.NotFound, "Device not found")
	case -1502:
		// Device is not active (disabled or deleted)
		return status.Error(codes.FailedPrecondition, "Device is not active")
	case -1503:
		// Device is already added to another trade point
		return status.Error(codes.AlreadyExists, "Device is already registered to another trade point")
	case -1601:
		// Purchase not found
		return status.Error(codes.NotFound, "Payment not found")
	case -14000002:
		// No trade points, need to create a trade point in the Kaspi Pay application
		return status.Error(codes.FailedPrecondition, "No trade points available. Please create a trade point in the Kaspi Pay application")
	case -99000002:
		// Trade point not found
		return status.Error(codes.NotFound, "Trade point not found")
	case -99000005:
		// Refund amount cannot exceed the purchase amount
		return status.Error(codes.InvalidArgument, "Refund amount cannot exceed the purchase amount")
	case -99000006:
		// Refund error, need to try again and contact the bank if the error persists
		return status.Error(codes.Internal, "Payment refund error. Please try again later")
	case 990000018:
		// Trade point disabled
		return status.Error(codes.FailedPrecondition, "Trade point is disabled")
	case 990000026:
		// Trade point does not accept payment with QR
		return status.Error(codes.FailedPrecondition, "Trade point does not accept QR payments")
	case 990000028:
		// Invalid operation amount specified
		return status.Error(codes.InvalidArgument, "Invalid payment amount")
	case 990000033:
		// No available payment methods
		return status.Error(codes.FailedPrecondition, "No available payment methods")
	case -99000001:
		// Purchase with the specified identifier not found
		return status.Error(codes.NotFound, "Payment with the specified ID not found")
	case -99000003:
		// The purchase trade point does not match the current device
		return status.Error(codes.PermissionDenied, "Payment trade point does not match current device")
	case -99000011:
		// Unable to return purchase (inappropriate purchase status)
		return status.Error(codes.FailedPrecondition, "Payment cannot be refunded due to its current status")
	case -99000020:
		// Partial refund not possible
		return status.Error(codes.FailedPrecondition, "Partial refund is not possible for this payment")
	case -999:
		// Service temporarily unavailable
		return status.Error(codes.Unavailable, "Kaspi Pay service is temporarily unavailable")
	case -10000:
		// No client certificate
		return status.Error(codes.Unauthenticated, "Authentication error: No client certificate")
	default:
		// Unknown error
		return status.Error(codes.Unknown, "Unexpected error from payment system: "+kaspiErr.Message)
	}
}
