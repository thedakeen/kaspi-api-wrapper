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

type KaspiService struct {
	log        *slog.Logger
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewKaspiService(log *slog.Logger, baseURL string, apiKey string) *KaspiService {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &KaspiService{
		log:        log,
		baseURL:    baseURL,
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

// generateRequestID generates X-Request-ID (2.1)
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// request makes a general request to the Kaspi API
func (s *KaspiService) request(ctx context.Context, method, path string, body, result any) error {
	const op = "service.kaspi.request"

	url := s.baseURL + path

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
	req.Header.Set("Api-Key", s.apiKey)

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
		return fmt.Errorf("%s:%w", op, err)
	}

	if baseResp.StatusCode != 0 {
		return fmt.Errorf("API error: %d - %s", baseResp.StatusCode, baseResp.Message)
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
