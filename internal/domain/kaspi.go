package domain

// BaseResponse represents base structure of response
type BaseResponse struct {
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message,omitempty"`
	Data       any    `json:"Data,omitempty"`
}

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
