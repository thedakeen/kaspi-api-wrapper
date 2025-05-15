package domain

type EnhancedDeviceRegisterRequest struct {
	DeviceID        string `json:"DeviceId"`
	TradePointID    int64  `json:"TradePointId"`
	OrganizationBin string `json:"OrganizationBin"`
}

type EnhancedDeviceDeleteRequest struct {
	DeviceToken     string `json:"DeviceToken"`
	OrganizationBin string `json:"OrganizationBin"`
}

type EnhancedQRCreateRequest struct {
	DeviceToken     string  `json:"DeviceToken"`
	Amount          float64 `json:"Amount"`
	ExternalID      string  `json:"ExternalId,omitempty"`
	OrganizationBin string  `json:"OrganizationBin"`
}

type EnhancedPaymentLinkCreateRequest struct {
	DeviceToken     string  `json:"DeviceToken"`
	Amount          float64 `json:"Amount"`
	ExternalID      string  `json:"ExternalId,omitempty"`
	OrganizationBin string  `json:"OrganizationBin"`
}

type EnhancedRefundRequest struct {
	DeviceToken     string  `json:"DeviceToken"`
	QrPaymentID     int64   `json:"QrPaymentId"`
	Amount          float64 `json:"Amount"`
	OrganizationBin string  `json:"OrganizationBin"`
}

type ClientInfoRequest struct {
	PhoneNumber string `json:"PhoneNumber"`
	DeviceToken int64  `json:"DeviceToken"`
}

type ClientInfoResponse struct {
	ClientName string `json:"ClientName"`
}

type RemotePaymentRequest struct {
	OrganizationBin string  `json:"OrganizationBin"`
	Amount          float64 `json:"Amount"`
	PhoneNumber     string  `json:"PhoneNumber"`
	DeviceToken     int64   `json:"DeviceToken"`
	Comment         string  `json:"Comment,omitempty"`
}

type RemotePaymentResponse struct {
	QrPaymentID int64 `json:"QrPaymentId"`
}

type RemotePaymentCancelRequest struct {
	OrganizationBin string `json:"OrganizationBin"`
	QrPaymentID     int64  `json:"QrPaymentId"`
	DeviceToken     int64  `json:"DeviceToken"`
}

type RemotePaymentCancelResponse struct {
	Status string `json:"Status"`
}
