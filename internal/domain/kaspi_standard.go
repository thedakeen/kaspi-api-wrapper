package domain

import "time"

//////// 	Refund domains	(standard)	////////

type QRRefundCreateRequest struct {
	DeviceToken string `json:"DeviceToken"`
	ExternalID  string `json:"ExternalId,omitempty"`
}

type QRRefundBehaviorOptions struct {
	QrCodeScanEventPollingInterval int `json:"QrCodeScanEventPollingInterval"`
	QrCodeScanWaitTimeout          int `json:"QrCodeScanWaitTimeout"`
}

type QRRefundCreateResponse struct {
	QrToken                 string                  `json:"QrToken"`
	ExpireDate              time.Time               `json:"ExpireDate"`
	QrReturnID              int64                   `json:"QrReturnId"`
	QrRefundBehaviorOptions QRRefundBehaviorOptions `json:"QrReturnBehaviorOptions"`
}

type RefundStatusResponse struct {
	Status string `json:"Status"`
}

type CustomerOperationsRequest struct {
	DeviceToken string `json:"DeviceToken"`
	QrReturnID  int64  `json:"QrReturnId"`
	MaxResult   int64  `json:"MaxResult,omitempty"`
}

type CustomerOperation struct {
	QrPaymentID     int64     `json:"QrPaymentId"`
	TransactionDate time.Time `json:"TransactionDate"`
	Amount          float64   `json:"Amount"`
}

type PaymentDetailsRequest struct {
	QrPaymentID int64  `json:"QrPaymentId"`
	DeviceToken string `json:"DeviceToken"`
}

type PaymentDetailsResponse struct {
	QrPaymentID           int64     `json:"QrPaymentId"`
	TotalAmount           float64   `json:"TotalAmount"`
	AvailableReturnAmount float64   `json:"AvailableReturnAmount"`
	TransactionDate       time.Time `json:"TransactionDate"`
}

type RefundRequest struct {
	DeviceToken string  `json:"DeviceToken"`
	QrPaymentID int64   `json:"QrPaymentId"`
	QrReturnID  int64   `json:"QrReturnId"`
	Amount      float64 `json:"Amount"`
}

type RefundResponse struct {
	ReturnOperationID int64 `json:"ReturnOperationId"`
}

//////// 	End of refund domains	(standard)	////////
