package domain

type EnhancedRefundRequest struct {
	DeviceToken     string  `json:"DeviceToken"`
	QrPaymentID     int64   `json:"QrPaymentId"`
	Amount          float64 `json:"Amount"`
	OrganizationBin string  `json:"OrganizationBin"`
}

type ClientInfoRequest struct {
	PhoneNumber string `json:"PhoneNumber"`
	DeviceToken string `json:"DeviceToken"`
}

type ClientInfoResponse struct {
	ClientName string `json:"ClientName"`
}

type RemotePaymentRequest struct {
	OrganizationBin string  `json:"OrganizationBin"`
	Amount          float64 `json:"Amount"`
	PhoneNumber     string  `json:"PhoneNumber"`
	DeviceToken     string  `json:"DeviceToken"`
	Comment         string  `json:"Comment,omitempty"`
}

type RemotePaymentResponse struct {
	QrPaymentID int64 `json:"QrPaymentId"`
}

type RemotePaymentCancelRequest struct {
	OrganizationBin string `json:"OrganizationBin"`
	QrPaymentID     int64  `json:"QrPaymentId"`
	DeviceToken     string `json:"DeviceToken"`
}

type RemotePaymentCancelResponse struct {
	Status string `json:"Status"`
}
