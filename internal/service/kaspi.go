package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kaspi-api-wrapper/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type KaspiService struct {
	log          *slog.Logger
	scheme       string
	baseURLBasic string
	baseURLStd   string
	baseURLEnh   string
	httpClient   HTTPClient
	apiKey       string
}

func NewKaspiService(log *slog.Logger,
	scheme string,
	baseURLBasic string,
	baseURLStd string,
	baseURLEnh string,
	apiKey string,
) *KaspiService {
	var httpClient *http.Client

	switch scheme {
	case "basic":
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	case "standard", "enhanced":
		//cert, err := tls.LoadX509KeyPair(certFile, certFile)
		//if err != nil {
		//	log.Error("failed to load certificate", "error", err)
		//	httpClient = &http.Client{Timeout: 30 * time.Second}
		//} else {
		//	tlsConfig := &tls.Config{
		//		Certificates: []tls.Certificate{cert},
		//	}
		//	transport := &http.Transport{TLSClientConfig: tlsConfig}
		//	httpClient = &http.Client{
		//		Transport: transport,
		//		Timeout:   30 * time.Second,
		//	}
		//}
	default:
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &KaspiService{
		log:          log,
		scheme:       scheme,
		baseURLBasic: baseURLBasic,
		baseURLStd:   baseURLStd,
		baseURLEnh:   baseURLEnh,
		httpClient:   httpClient,
		apiKey:       apiKey,
	}
}

// generateRequestID generates X-Request-ID (2.1)
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// SetHTTPClient sets the HTTP client for testing
func (s *KaspiService) SetHTTPClient(client HTTPClient) {
	s.httpClient = client
}

// GetBaseURL retrieves the base URL based on the current scheme
func (s *KaspiService) GetBaseURL() string {
	switch s.scheme {
	case "basic":
		return s.baseURLBasic
	case "standard":
		return s.baseURLStd
	case "enhanced":
		return s.baseURLEnh
	default:
		return s.baseURLBasic
	}
}

// Request makes a general request to the Kaspi API, exposed method for testing
func (s *KaspiService) Request(ctx context.Context, method, path string, body, result any) error {
	const op = "service.kaspi.request"

	url := s.GetBaseURL() + path

	log := s.log.With(
		slog.String("op", op),
		slog.String("method", method),
		slog.String("url", url),
	)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", generateRequestID())

	// Api-Key for request via first scheme
	if s.scheme == "basic" {
		req.Header.Set("Api-Key", s.apiKey)
	}

	log.Debug("sending request")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Debug("received response", "status", resp.Status, "body", string(respBody))

	var baseResp domain.BaseResponse
	err = json.Unmarshal(respBody, &baseResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if baseResp.StatusCode != 0 {
		return &domain.KaspiError{
			StatusCode: baseResp.StatusCode,
			Message:    baseResp.Message,
		}
	}

	if result != nil && baseResp.Data != nil {
		dataJSON, err := json.Marshal(baseResp.Data)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}

		err = json.Unmarshal(dataJSON, result)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
	}

	return nil
}

// request makes a general request to the Kaspi API
func (s *KaspiService) request(ctx context.Context, method, path string, body, result any) error {
	return s.Request(ctx, method, path, body, result)
}

//////// 	Device service methods	////////

// GetTradePoints retrieves list of trade points from Kaspi API (2.2.2)
func (s *KaspiService) GetTradePoints(ctx context.Context) ([]domain.TradePoint, error) {
	const op = "service.kaspi.GetTradePoints"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("getting all trade points")

	path := "/partner/tradepoints"

	var result []domain.TradePoint
	err := s.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug("all trade points got successfully")

	return result, nil
}

// RegisterDevice registers a new device in Kaspi Pay (2.2.3)
func (s *KaspiService) RegisterDevice(ctx context.Context, req domain.DeviceRegisterRequest) (*domain.DeviceRegisterResponse, error) {
	const op = "service.kaspi.RegisterDevice"

	log := s.log.With(
		slog.String("op", op),
		slog.String("Device", req.DeviceID),
		slog.Int64("Trade Point", req.TradePointID),
	)

	log.Debug("registering new device")

	path := "/device/register"

	var result domain.DeviceRegisterResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug("new device registered successfully")

	return &result, nil
}

// DeleteDevice deletes a device from Kaspi Pay (2.2.4)
func (s *KaspiService) DeleteDevice(ctx context.Context, deviceToken string) error {
	const op = "service.kaspi.DeleteDevice"

	log := s.log.With(
		slog.String("op", op),
		slog.String("DeviceToken", deviceToken),
	)

	log.Debug("deleting device")

	path := "/device/delete"

	req := domain.DeviceRegisterResponse{
		DeviceToken: deviceToken,
	}

	err := s.request(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Debug("device deleted successfully")

	return nil
}

//////// 	End of device service methods	////////

//////// 	Payment service	methods	////////

func (s *KaspiService) CreateQR(ctx context.Context, req domain.QRCreateRequest) (*domain.QRCreateResponse, error) {
	const op = "service.kaspi.CreateQR"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
	)

	log.Debug("creating QR token for payment")

	path := "/qr/create"

	var result domain.QRCreateResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("QR token created successfully")

	return &result, nil
}

// CreatePaymentLink creates a payment link (2.3.2)
func (s *KaspiService) CreatePaymentLink(ctx context.Context, req domain.PaymentLinkCreateRequest) (*domain.PaymentLinkCreateResponse, error) {
	const op = "service.kaspi.CreatePaymentLink"

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
	)

	log.Debug("creating payment link")

	path := "/qr/create-link"

	var result domain.PaymentLinkCreateResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment link created successfully")

	return &result, nil
}

// GetPaymentStatus retrieves the status of a payment (2.3.3)
func (s *KaspiService) GetPaymentStatus(ctx context.Context, qrPaymentID int64) (*domain.PaymentStatusResponse, error) {
	const op = "service.kaspi.GetPaymentStatus"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("qrPaymentID", qrPaymentID),
	)

	log.Debug("getting payment status")

	path := fmt.Sprintf("/payment/status/%d", qrPaymentID)

	var result domain.PaymentStatusResponse
	err := s.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment status retrieved successfully", "status", result.Status)

	return &result, nil
}

//////// 	End of payment service	methods	////////
