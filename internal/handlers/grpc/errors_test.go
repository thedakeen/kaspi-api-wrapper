package grpchandler_test

import (
	"errors"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"kaspi-api-wrapper/internal/domain"
	grpchandler "kaspi-api-wrapper/internal/handlers/grpc"
	"kaspi-api-wrapper/internal/validator"
	"log/slog"
	"os"
)

func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func TestHandleError(t *testing.T) {
	log := setupTestLogger()

	t.Run("handles unsupported feature error", func(t *testing.T) {
		err := domain.ErrUnsupportedFeature

		result := grpchandler.HandleError(err, log)

		st, ok := status.FromError(result)
		if !ok {
			t.Fatal("Expected gRPC status error")
		}

		if st.Code() != codes.PermissionDenied {
			t.Errorf("Expected code PermissionDenied, got %s", st.Code())
		}
	})

	t.Run("handles validation error", func(t *testing.T) {
		err := &validator.ValidationError{
			Field:   "deviceId",
			Message: "device ID is required",
			Err:     validator.ErrRequiredField,
		}

		result := grpchandler.HandleError(err, log)

		st, ok := status.FromError(result)
		if !ok {
			t.Fatal("Expected gRPC status error")
		}

		if st.Code() != codes.InvalidArgument {
			t.Errorf("Expected code InvalidArgument, got %s", st.Code())
		}

		expectedMsg := "deviceId: device ID is required"
		if !strings.Contains(st.Message(), expectedMsg) {
			t.Errorf("Expected message to contain '%s', got '%s'", expectedMsg, st.Message())
		}
	})

	testCases := []struct {
		name         string
		err          error
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name:         "Device not found",
			err:          &domain.KaspiError{StatusCode: -1501, Message: "Device not found"},
			expectedCode: codes.NotFound,
			expectedMsg:  "Device not found",
		},
		{
			name:         "Device not active",
			err:          &domain.KaspiError{StatusCode: -1502, Message: "Device not active"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "Device is not active",
		},
		{
			name:         "Device already registered",
			err:          &domain.KaspiError{StatusCode: -1503, Message: "Device registered elsewhere"},
			expectedCode: codes.AlreadyExists,
			expectedMsg:  "Device is already registered to another trade point",
		},
		{
			name:         "Payment not found",
			err:          &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"},
			expectedCode: codes.NotFound,
			expectedMsg:  "Payment not found",
		},
		{
			name:         "No trade points",
			err:          &domain.KaspiError{StatusCode: -14000002, Message: "No trade points"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "No trade points available",
		},
		{
			name:         "Trade point not found",
			err:          &domain.KaspiError{StatusCode: -99000002, Message: "Trade point not found"},
			expectedCode: codes.NotFound,
			expectedMsg:  "Trade point not found",
		},
		{
			name:         "Refund amount exceeds purchase",
			err:          &domain.KaspiError{StatusCode: -99000005, Message: "Refund amount too high"},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "Refund amount cannot exceed the purchase amount",
		},
		{
			name:         "Refund error",
			err:          &domain.KaspiError{StatusCode: -99000006, Message: "Refund failed"},
			expectedCode: codes.Internal,
			expectedMsg:  "Payment refund error",
		},
		{
			name:         "Trade point disabled",
			err:          &domain.KaspiError{StatusCode: 990000018, Message: "Trade point disabled"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "Trade point is disabled",
		},
		{
			name:         "No QR payments",
			err:          &domain.KaspiError{StatusCode: 990000026, Message: "No QR payments"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "Trade point does not accept QR payments",
		},
		{
			name:         "Invalid amount",
			err:          &domain.KaspiError{StatusCode: 990000028, Message: "Invalid amount"},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "Invalid payment amount",
		},
		{
			name:         "No payment methods",
			err:          &domain.KaspiError{StatusCode: 990000033, Message: "No payment methods"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "No available payment methods",
		},
		{
			name:         "Payment ID not found",
			err:          &domain.KaspiError{StatusCode: -99000001, Message: "Payment ID not found"},
			expectedCode: codes.NotFound,
			expectedMsg:  "Payment with the specified ID not found",
		},
		{
			name:         "Trade point mismatch",
			err:          &domain.KaspiError{StatusCode: -99000003, Message: "Trade point mismatch"},
			expectedCode: codes.PermissionDenied,
			expectedMsg:  "Payment trade point does not match current device",
		},
		{
			name:         "Cannot refund due to status",
			err:          &domain.KaspiError{StatusCode: -99000011, Message: "Cannot refund"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "Payment cannot be refunded due to its current status",
		},
		{
			name:         "Partial refund not possible",
			err:          &domain.KaspiError{StatusCode: -99000020, Message: "No partial refund"},
			expectedCode: codes.FailedPrecondition,
			expectedMsg:  "Partial refund is not possible for this payment",
		},
		{
			name:         "Service unavailable",
			err:          &domain.KaspiError{StatusCode: -999, Message: "Service unavailable"},
			expectedCode: codes.Unavailable,
			expectedMsg:  "Kaspi Pay service is temporarily unavailable",
		},
		{
			name:         "No client certificate",
			err:          &domain.KaspiError{StatusCode: -10000, Message: "No client certificate"},
			expectedCode: codes.Unauthenticated,
			expectedMsg:  "Authentication error: No client certificate",
		},
		{
			name:         "Unknown Kaspi error",
			err:          &domain.KaspiError{StatusCode: -99999, Message: "Unknown error"},
			expectedCode: codes.Unknown,
			expectedMsg:  "Unexpected error from payment system: Unknown error",
		},
		{
			name:         "Non-Kaspi error",
			err:          errors.New("Some other error"),
			expectedCode: codes.Internal,
			expectedMsg:  "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grpchandler.HandleError(tc.err, log)

			st, ok := status.FromError(result)
			if !ok {
				t.Fatal("Expected gRPC status error")
			}

			if st.Code() != tc.expectedCode {
				t.Errorf("Expected code %s, got %s", tc.expectedCode, st.Code())
			}

			if !strings.Contains(st.Message(), tc.expectedMsg) {
				t.Errorf("Expected message to contain '%s', got '%s'", tc.expectedMsg, st.Message())
			}
		})
	}
}
