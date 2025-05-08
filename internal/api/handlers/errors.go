package handlers

import (
	"errors"
	"kaspi-api-wrapper/internal/domain"
	"log/slog"
	"net/http"
)

// HandleKaspiError handles Kaspi API errors and maps them to appropriate HTTP responses
func HandleKaspiError(w http.ResponseWriter, err error, log *slog.Logger) {
	if err != nil && errors.Is(err, domain.ErrUnsupportedFeature) {
		log.Error("scheme compatibility error", "error", err)
		ForbiddenError(w, err.Error())
		return
	}
	kaspiErr, ok := domain.IsKaspiError(err)
	if !ok {
		log.Error("unexpected error", "error", err)
		InternalServerError(w, "Internal server error")
		return
	}

	log.Error("kaspi API error",
		"status_code", kaspiErr.StatusCode,
		"message", kaspiErr.Message)

	switch kaspiErr.StatusCode {
	case -1501:
		// Device with the specified identifier not found
		NotFoundError(w, "Device not found")
	case -1502:
		// Device is not active (disabled or deleted)
		BadRequestError(w, "Device is not active")
	case -1503:
		// Device is already added to another trade point
		ConflictError(w, "Device is already registered to another trade point")
	case -1601:
		// Purchase not found
		NotFoundError(w, "Payment not found")
	case -14000002:
		// No trade points, need to create a trade point in the Kaspi Pay application
		BadRequestError(w, "No trade points available. Please create a trade point in the Kaspi Pay application")
	case -99000002:
		// Trade point not found
		NotFoundError(w, "Trade point not found")
	case -99000005:
		// Refund amount cannot exceed the purchase amount
		BadRequestError(w, "Refund amount cannot exceed the purchase amount")
	case -99000006:
		// Refund error, need to try again and contact the bank if the error persists
		InternalServerError(w, "Payment refund error. Please try again later")
	case 990000018:
		// Trade point disabled
		BadRequestError(w, "Trade point is disabled")
	case 990000026:
		// Trade point does not accept payment with QR
		BadRequestError(w, "Trade point does not accept QR payments")
	case 990000028:
		// Invalid operation amount specified
		BadRequestError(w, "Invalid payment amount")
	case 990000033:
		// No available payment methods
		BadRequestError(w, "No available payment methods")
	case -99000001:
		// Purchase with the specified identifier not found
		NotFoundError(w, "Payment with the specified ID not found")
	case -99000003:
		// The purchase trade point does not match the current device
		ForbiddenError(w, "Payment trade point does not match current device")
	case -99000011:
		// Unable to return purchase (inappropriate purchase status)
		BadRequestError(w, "Payment cannot be refunded due to its current status")
	case -99000020:
		// Partial refund not possible
		BadRequestError(w, "Partial refund is not possible for this payment")
	case -999:
		// Service temporarily unavailable
		ServiceUnavailableError(w, "Kaspi Pay service is temporarily unavailable")
	case -10000:
		// No client certificate
		UnauthorizedError(w, "Authentication error: No client certificate")
	default:
		// Unknown error
		InternalServerError(w, "Unexpected error from payment system: "+kaspiErr.Message)
	}
}
