// internal/handlers/grpc/middleware/scheme.go
package middleware

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var methodRequirements = map[string]string{
	// Basic scheme methods (1)
	"/kaspi.api.v1.DeviceService/GetTradePoints":      "basic",
	"/kaspi.api.v1.DeviceService/RegisterDevice":      "basic",
	"/kaspi.api.v1.DeviceService/DeleteDevice":        "basic",
	"/kaspi.api.v1.PaymentService/CreateQR":           "basic",
	"/kaspi.api.v1.PaymentService/CreatePaymentLink":  "basic",
	"/kaspi.api.v1.PaymentService/GetPaymentStatus":   "basic",
	"/kaspi.api.v1.UtilityService/HealthCheck":        "basic",
	"/kaspi.api.v1.UtilityService/TestScanQR":         "basic",
	"/kaspi.api.v1.UtilityService/TestConfirmPayment": "basic",
	"/kaspi.api.v1.UtilityService/TestScanError":      "basic",
	"/kaspi.api.v1.UtilityService/TestConfirmError":   "basic",

	// Standard scheme methods (2)
	"/kaspi.api.v1.RefundService/CreateRefundQR":        "standard",
	"/kaspi.api.v1.RefundService/GetRefundStatus":       "standard",
	"/kaspi.api.v1.RefundService/GetCustomerOperations": "standard",
	"/kaspi.api.v1.RefundService/GetPaymentDetails":     "standard",
	"/kaspi.api.v1.RefundService/RefundPayment":         "standard",

	// Enhanced scheme methods (3)
	"/kaspi.api.v1.DeviceService/GetTradePointsEnhanced":        "enhanced",
	"/kaspi.api.v1.DeviceService/RegisterDeviceEnhanced":        "enhanced",
	"/kaspi.api.v1.DeviceService/DeleteDeviceEnhanced":          "enhanced",
	"/kaspi.api.v1.PaymentService/CreateQREnhanced":             "enhanced",
	"/kaspi.api.v1.PaymentService/CreatePaymentLinkEnhanced":    "enhanced",
	"/kaspi.api.v1.EnhancedRefundService/RefundPaymentEnhanced": "enhanced",
	"/kaspi.api.v1.EnhancedRefundService/GetClientInfo":         "enhanced",
	"/kaspi.api.v1.EnhancedRefundService/CreateRemotePayment":   "enhanced",
	"/kaspi.api.v1.EnhancedRefundService/CancelRemotePayment":   "enhanced",
}

// isMethodAllowed checks if the method is allowed in the current scheme
func isMethodAllowed(fullMethod string, currentScheme string) bool {
	requiredScheme, exists := methodRequirements[fullMethod]
	if !exists {
		return false
	}

	switch requiredScheme {
	case "basic":
		return true
	case "standard":
		return currentScheme == "standard" || currentScheme == "enhanced"
	case "enhanced":
		return currentScheme == "enhanced"
	default:
		return false
	}
}

// SchemeInterceptor creates a gRPC interceptor that restricts access based on scheme level
func SchemeInterceptor(currentScheme string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !isMethodAllowed(info.FullMethod, currentScheme) {
			methodName := strings.Split(info.FullMethod, "/")
			shortName := methodName[len(methodName)-1]

			message := fmt.Sprintf("Method %s requires %s scheme, but current scheme is %s",
				shortName, methodRequirements[info.FullMethod], currentScheme)

			return nil, status.Error(codes.PermissionDenied, message)
		}

		return handler(ctx, req)
	}
}
