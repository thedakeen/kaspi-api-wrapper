package service

import (
	"context"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	"kaspi-api-wrapper/internal/validator"
	"log/slog"
	"net/http"
)

// HealthCheck checks the availability of the Kaspi API (5.1)
func (s *KaspiService) HealthCheck(ctx context.Context) error {
	const op = "service.kaspi.HealthCheck"

	log := s.log.With(slog.String("op", op))
	log.Debug("checking Kaspi API health")

	path := "/health/ping"

	err := s.request(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("Kaspi API health check successful")
	return nil
}

// TestScanQR simulates scanning a QR code (5.2)
func (s *KaspiService) TestScanQR(ctx context.Context, req domain.TestScanRequest) error {
	const op = "service.kaspi.TestScanQR"

	log := s.log.With(
		slog.String("op", op),
		slog.String("qrPaymentId", req.QrPaymentID),
	)

	if err := validator.ValidateTestScanRequest(req); err != nil {
		log.Warn("invalid test scan qr request", "error", err.Error())
		return err
	}

	log.Debug("simulating QR code scan")

	path := "/test/payment/scan"

	err := s.request(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("QR code scan simulation successful")
	return nil
}

// TestConfirmPayment simulates payment confirmation (5.3)
func (s *KaspiService) TestConfirmPayment(ctx context.Context, req domain.TestConfirmRequest) error {
	const op = "service.kaspi.TestConfirmPayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("qrPaymentId", req.QrPaymentID),
	)

	if err := validator.ValidateTestConfirmRequest(req); err != nil {
		log.Warn("invalid test confirm request", "error", err.Error())
		return err
	}

	log.Debug("simulating payment confirmation")

	path := "/test/payment/confirm"

	err := s.request(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment confirmation simulation successful")
	return nil
}

// TestScanError simulates an error during QR code scanning (5.4)
func (s *KaspiService) TestScanError(ctx context.Context, req domain.TestScanErrorRequest) error {
	const op = "service.kaspi.TestScanError"

	log := s.log.With(
		slog.String("op", op),
		slog.String("qrPaymentId", req.QrPaymentID),
	)

	if err := validator.ValidateTestScanErrorRequest(req); err != nil {
		log.Warn("invalid test scan error request", "error", err.Error())
		return err
	}

	log.Debug("simulating QR code scan error")

	path := "/test/payment/scanerror"

	err := s.request(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("QR code scan error simulation successful")
	return nil
}

// TestConfirmError simulates an error during payment confirmation (5.5)
func (s *KaspiService) TestConfirmError(ctx context.Context, req domain.TestConfirmErrorRequest) error {
	const op = "service.kaspi.TestConfirmError"

	log := s.log.With(
		slog.String("op", op),
		slog.String("qrPaymentId", req.QrPaymentID),
	)

	if err := validator.ValidateTestConfirmErrorRequest(req); err != nil {
		log.Warn("invalid test confirm error request", "error", err.Error())
		return err
	}

	log.Debug("simulating payment confirmation error")

	path := "/test/payment/confirmerror"

	err := s.request(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("payment confirmation error simulation successful")
	return nil
}
