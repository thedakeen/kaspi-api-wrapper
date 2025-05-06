package domain

type TestScanRequest struct {
	QrPaymentID string `json:"qrPaymentId"`
}

type TestConfirmRequest struct {
	QrPaymentID string `json:"qrPaymentId"`
}

type TestScanErrorRequest struct {
	QrPaymentID string `json:"qrPaymentId"`
}

type TestConfirmErrorRequest struct {
	QrPaymentID string `json:"qrPaymentId"`
}
