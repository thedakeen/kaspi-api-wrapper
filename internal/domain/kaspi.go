package domain

import "time"

// BaseResponse represents base structure of response
type BaseResponse struct {
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message,omitempty"`
	Data       any    `json:"Data,omitempty"`
}

//////// 	Device domains		////////

type TradePoint struct {
	TradePointID   int64  `json:"TradePointId"`
	TradePointName string `json:"TradePointName"`
}

type DeviceRegisterRequest struct {
	DeviceID     string `json:"DeviceId"`
	TradePointID int64  `json:"TradePointId"`
}

type DeviceRegisterResponse struct {
	DeviceToken string `json:"DeviceToken"`
}

//////// 	End of device domains		////////

//////// 	Payment domains		////////

type QRCreateRequest struct {
	DeviceToken string  `json:"DeviceToken"`
	Amount      float64 `json:"Amount"`
	ExternalID  string  `json:"ExternalId,omitempty"`
}

type QRPaymentBehaviorOptions struct {
	StatusPollingInterval      int `json:"StatusPollingInterval"`
	QrCodeScanWaitTimeout      int `json:"QrCodeScanWaitTimeout"`
	PaymentConfirmationTimeout int `json:"PaymentConfirmationTimeout"`
}

type QRCreateResponse struct {
	QrToken                  string                   `json:"QrToken"`
	ExpireDate               time.Time                `json:"ExpireDate"`
	QrPaymentID              int64                    `json:"QrPaymentId"`
	PaymentMethods           []string                 `json:"PaymentMethods"`
	QrPaymentBehaviorOptions QRPaymentBehaviorOptions `json:"QrPaymentBehaviorOptions"`
}

type PaymentLinkCreateRequest struct {
	DeviceToken string  `json:"DeviceToken"`
	Amount      float64 `json:"Amount"`
	ExternalID  string  `json:"ExternalId,omitempty"`
}

type PaymentBehaviorOptions struct {
	StatusPollingInterval      int `json:"StatusPollingInterval"`
	LinkActivationWaitTimeout  int `json:"LinkActivationWaitTimeout"`
	PaymentConfirmationTimeout int `json:"PaymentConfirmationTimeout"`
}

type PaymentLinkCreateResponse struct {
	PaymentLink            string                 `json:"PaymentLink"`
	ExpireDate             time.Time              `json:"ExpireDate"`
	PaymentID              int64                  `json:"PaymentId"`
	PaymentMethods         []string               `json:"PaymentMethods"`
	PaymentBehaviorOptions PaymentBehaviorOptions `json:"PaymentBehaviorOptions"`
}

type PaymentStatusResponse struct {
	Status        string  `json:"Status"`
	TransactionID string  `json:"TransactionId,omitempty"`
	LoanOfferName string  `json:"LoanOfferName,omitempty"`
	LoanTerm      int     `json:"LoanTerm,omitempty"`
	IsOffer       bool    `json:"IsOffer,omitempty"`
	ProductType   string  `json:"ProductType,omitempty"`
	Amount        float64 `json:"Amount,omitempty"`
	StoreName     string  `json:"StoreName,omitempty"`
	Address       string  `json:"Address,omitempty"`
	City          string  `json:"City,omitempty"`
}

//////// 	End of payment domains		////////
