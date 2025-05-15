package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/storage"
	"kaspi-api-wrapper/internal/validator"
	"log/slog"
	"net/http"
	"os"
	"software.sslmate.com/src/go-pkcs12"
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

	deviceSaver DeviceSaver
}

// TLSConfig for scheme 2 & 3
type TLSConfig struct {
	PfxFile       string // .pfx file containing both certificate and private key
	Password      string // password for the private key
	RootCAFile    string // root CA certificate
	UseClientCert bool   // whether to use client certificate authentication
}

type DeviceSaver interface {
	SaveDevice(ctx context.Context, deviceID string, deviceToken string, tradePointID int64) error
	SaveDeviceEnhanced(ctx context.Context, deviceID string, deviceToken string, tradePointID int64, organizationBin string) error
}

func NewKaspiService(log *slog.Logger,
	scheme string,
	baseURLBasic string,
	baseURLStd string,
	baseURLEnh string,
	apiKey string,
	tlsConfig *TLSConfig,

	deviceSaver DeviceSaver,
) *KaspiService {
	var httpClient *http.Client
	var err error

	switch scheme {
	case "basic":
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	case "standard", "enhanced":
		if tlsConfig == nil {
			log.Error("TLS config not provided for scheme requiring client certificates",
				"scheme", scheme)
			httpClient = &http.Client{Timeout: 30 * time.Second}
		} else {
			tlsConfig.UseClientCert = true
			httpClient, err = loadTLSConfig(log, tlsConfig)
			if err != nil {
				panic(err)
			}

		}
		//tlsConfig, err := loadTLSConfig(scheme)
		//if err != nil {
		//	log.Error("failed to load TLS config", "error", err)
		//	httpClient = &http.Client{Timeout: 30 * time.Second}
		//} else {
		//	transport := &http.Transport{
		//		TLSClientConfig: tlsConfig,
		//	}
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

		deviceSaver: deviceSaver,
	}
}

func loadTLSConfig(log *slog.Logger, cfg *TLSConfig) (*http.Client, error) {
	const op = "service.kaspi.loadTLSConfig"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	if !cfg.UseClientCert {
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: tr,
		}, nil
	}

	var clientCert tls.Certificate

	if cfg.PfxFile != "" {
		pfxData, err := os.ReadFile(cfg.PfxFile) // ioutil
		if err != nil {
			return nil, fmt.Errorf("%s: failed to read PFX file: %w", op, err)
		}

		privateKey, certificate, caCerts, err := pkcs12.DecodeChain(pfxData, cfg.Password)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to parse PFX data: %w", op, err)
		}

		clientCert.Certificate = make([][]byte, len(caCerts)+1)
		clientCert.Certificate[0] = certificate.Raw
		for i, ca := range caCerts {
			clientCert.Certificate[i+1] = ca.Raw
		}
		clientCert.PrivateKey = privateKey
	} else {
		return nil, fmt.Errorf("%s: no valid certificate configuration provided", op)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	if cfg.RootCAFile != "" {
		rootCA, err := os.ReadFile(cfg.RootCAFile) // ioutil
		if err != nil {
			return nil, fmt.Errorf("%s: failed to read root CA file: %w", op, err)
		}

		rootCAPool := x509.NewCertPool()
		if !rootCAPool.AppendCertsFromPEM(rootCA) {
			return nil, fmt.Errorf("%s: failed to append root CA to cert pool", op)
		}

		tlsConfig.RootCAs = rootCAPool
	}

	tr.TLSClientConfig = tlsConfig

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}, nil
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
	} else {
		log.Debug("using certificate auth", "scheme", s.scheme)
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
		fmt.Println(respBody)
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if baseResp.StatusCode != 0 {
		kaspiErr := &domain.KaspiError{
			StatusCode: baseResp.StatusCode,
			Message:    baseResp.Message,
		}

		if kaspiErr.StatusCode == -10000 {
			log.Error("certificate authentication failed - check your client certificate setup")
		}

		return kaspiErr
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

	if s.scheme == "enhanced" {
		return nil, domain.ErrUnsupportedFeature
	}

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

	if s.scheme == "enhanced" {
		return nil, domain.ErrUnsupportedFeature
	}

	log := s.log.With(
		slog.String("op", op),
		slog.String("Device", req.DeviceID),
		slog.Int64("Trade Point", req.TradePointID),
	)

	if err := validator.ValidateDeviceRegisterRequest(req); err != nil {
		log.Warn("invalid device register request", "error", err.Error())
		return nil, err
	}

	log.Debug("registering new device")

	path := "/device/register"

	var result domain.DeviceRegisterResponse
	err := s.request(ctx, http.MethodPost, path, req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug("new device registered successfully")

	// DB interaction
	log.Debug("saving device to database")

	err = s.deviceSaver.SaveDevice(ctx, req.DeviceID, result.DeviceToken, req.TradePointID)
	if err != nil {
		log.Error("failed to save device to database")
		switch {
		case errors.Is(err, storage.ErrDeviceExists):
			return nil, &domain.KaspiError{
				StatusCode: -1503,
				Message:    "Device is already added to another trade point",
			}
		default:
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Debug("device saved to database successfully")

	return &result, nil
}

// DeleteDevice deletes a device from Kaspi Pay (2.2.4)
func (s *KaspiService) DeleteDevice(ctx context.Context, deviceToken string) error {
	const op = "service.kaspi.DeleteDevice"

	if s.scheme == "enhanced" {
		return domain.ErrUnsupportedFeature
	}

	log := s.log.With(
		slog.String("op", op),
		slog.String("DeviceToken", deviceToken),
	)

	if err := validator.ValidateDeviceToken(deviceToken); err != nil {
		log.Warn("invalid device token", "error", err.Error())
		return err
	}

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

	if s.scheme == "enhanced" {
		return nil, domain.ErrUnsupportedFeature
	}

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
	)

	if err := validator.ValidateQRCreateRequest(req); err != nil {
		log.Warn("invalid QR create request", "error", err.Error())
		return nil, err
	}

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

	if s.scheme == "enhanced" {
		return nil, domain.ErrUnsupportedFeature
	}

	log := s.log.With(
		slog.String("op", op),
		slog.String("deviceToken", req.DeviceToken),
		slog.Float64("amount", req.Amount),
	)

	if err := validator.ValidatePaymentLinkCreateRequest(req); err != nil {
		log.Warn("invalid payment link create request", "error", err.Error())
		return nil, err
	}

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

	if qrPaymentID <= 0 {
		return nil, &validator.ValidationError{
			Field:   "qrPaymentId",
			Message: "Invalid payment ID format",
			Err:     validator.ErrInvalidID,
		}
	}

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
