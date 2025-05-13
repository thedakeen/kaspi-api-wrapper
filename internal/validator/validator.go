package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
)

// Validation errors
var (
	ErrRequiredField = fmt.Errorf("required field is missing")
	ErrInvalidValue  = fmt.Errorf("invalid value")
	ErrInvalidAmount = fmt.Errorf("invalid amount")
	ErrInvalidToken  = fmt.Errorf("invalid token")
	ErrInvalidID     = fmt.Errorf("invalid ID")
	ErrInvalidOrgBin = fmt.Errorf("invalid organization BIN")
	ErrInvalidPhone  = fmt.Errorf("invalid phone number")
)

// ValidationError represents a validation error with a field and message
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

// GRPCError handles validation errors in gRPC
func GRPCError(err error) error {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		return status.Error(codes.InvalidArgument, valErr.Error())
	}
	return status.Error(codes.Internal, "validation error")
}

// HTTPError handles validation errors in HTTP
func HTTPError(w http.ResponseWriter, err error) bool {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   valErr.Error(),
		})
		return true
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   "validation error",
	})
	return true
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateDeviceRegisterRequest validates a device registration request
func ValidateDeviceRegisterRequest(req domain.DeviceRegisterRequest) error {
	if req.DeviceID == "" {
		return &ValidationError{
			Field:   "deviceId",
			Message: "device ID is required",
			Err:     ErrRequiredField,
		}
	}

	if req.TradePointID <= 0 {
		return &ValidationError{
			Field:   "tradePointId",
			Message: "trade point ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	return nil
}

// ValidateDeviceToken validates a device token
func ValidateDeviceToken(deviceToken string) error {
	if deviceToken == "" {
		return &ValidationError{
			Field:   "deviceToken",
			Message: "device token is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateEnhancedDeviceRegisterRequest validates an enhanced device registration request
func ValidateEnhancedDeviceRegisterRequest(req domain.EnhancedDeviceRegisterRequest) error {
	if req.DeviceID == "" {
		return &ValidationError{
			Field:   "deviceId",
			Message: "device ID is required",
			Err:     ErrRequiredField,
		}
	}

	if req.TradePointID <= 0 {
		return &ValidationError{
			Field:   "tradePointId",
			Message: "trade point ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateEnhancedDeviceDeleteRequest validates an enhanced device deletion request
func ValidateEnhancedDeviceDeleteRequest(req domain.EnhancedDeviceDeleteRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// Payment Validation Functions

// ValidateQRCreateRequest validates a QR creation request
func ValidateQRCreateRequest(req domain.QRCreateRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	return nil
}

// ValidatePaymentLinkCreateRequest validates a payment link creation request
func ValidatePaymentLinkCreateRequest(req domain.PaymentLinkCreateRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	return nil
}

// ValidateEnhancedQRCreateRequest validates an enhanced QR creation request
func ValidateEnhancedQRCreateRequest(req domain.EnhancedQRCreateRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateEnhancedPaymentLinkCreateRequest validates an enhanced payment link creation request
func ValidateEnhancedPaymentLinkCreateRequest(req domain.EnhancedPaymentLinkCreateRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// Refund Validation Functions

// ValidateQRRefundCreateRequest validates a QR refund creation request
func ValidateQRRefundCreateRequest(req domain.QRRefundCreateRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	return nil
}

// ValidateCustomerOperationsRequest validates a customer operations request
func ValidateCustomerOperationsRequest(req domain.CustomerOperationsRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.QrReturnID <= 0 {
		return &ValidationError{
			Field:   "qrReturnId",
			Message: "QR return ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	return nil
}

// ValidatePaymentDetailsRequest validates a payment details request
func ValidatePaymentDetailsRequest(qrPaymentID int64, deviceToken string) error {
	if err := ValidateDeviceToken(deviceToken); err != nil {
		return err
	}

	if qrPaymentID <= 0 {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	return nil
}

// ValidateRefundRequest validates a refund request
func ValidateRefundRequest(req domain.RefundRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.QrPaymentID <= 0 {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	if req.QrReturnID <= 0 {
		return &ValidationError{
			Field:   "qrReturnId",
			Message: "QR return ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	return nil
}

// ValidateEnhancedRefundRequest validates an enhanced refund request
func ValidateEnhancedRefundRequest(req domain.EnhancedRefundRequest) error {
	if err := ValidateDeviceToken(req.DeviceToken); err != nil {
		return err
	}

	if req.QrPaymentID <= 0 {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateRemotePaymentRequest validates a remote payment request
func ValidateRemotePaymentRequest(req domain.RemotePaymentRequest) error {
	if req.DeviceToken <= 0 {
		return &ValidationError{
			Field:   "deviceToken",
			Message: "device token must be a positive number",
			Err:     ErrInvalidToken,
		}
	}

	if req.PhoneNumber == "" {
		return &ValidationError{
			Field:   "phoneNumber",
			Message: "phone number is required",
			Err:     ErrInvalidPhone,
		}
	}

	if req.Amount <= 0 {
		return &ValidationError{
			Field:   "amount",
			Message: "amount must be greater than zero",
			Err:     ErrInvalidAmount,
		}
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateRemotePaymentCancelRequest validates a remote payment cancel request
func ValidateRemotePaymentCancelRequest(req domain.RemotePaymentCancelRequest) error {
	if req.DeviceToken <= 0 {
		return &ValidationError{
			Field:   "deviceToken",
			Message: "device token must be a positive number",
			Err:     ErrInvalidToken,
		}
	}

	if req.QrPaymentID <= 0 {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID must be a positive number",
			Err:     ErrInvalidID,
		}
	}

	if req.OrganizationBin == "" {
		return &ValidationError{
			Field:   "organizationBin",
			Message: "organization BIN is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// Utility Validation Functions

// ValidateTestScanRequest validates a test scan request
func ValidateTestScanRequest(req domain.TestScanRequest) error {
	if req.QrPaymentID == "" {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateTestConfirmRequest validates a test confirm request
func ValidateTestConfirmRequest(req domain.TestConfirmRequest) error {
	if req.QrPaymentID == "" {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateTestScanErrorRequest validates a test scan error request
func ValidateTestScanErrorRequest(req domain.TestScanErrorRequest) error {
	if req.QrPaymentID == "" {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}

// ValidateTestConfirmErrorRequest validates a test confirm error request
func ValidateTestConfirmErrorRequest(req domain.TestConfirmErrorRequest) error {
	if req.QrPaymentID == "" {
		return &ValidationError{
			Field:   "qrPaymentId",
			Message: "QR payment ID is required",
			Err:     ErrRequiredField,
		}
	}

	return nil
}
