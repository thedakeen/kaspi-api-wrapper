package domain_test

import (
	"errors"
	"fmt"
	"kaspi-api-wrapper/internal/domain"
	"testing"
)

func TestKaspiError(t *testing.T) {
	t.Run("Error method returns formatted message", func(t *testing.T) {
		err := &domain.KaspiError{
			StatusCode: -1501,
			Message:    "Device not found",
		}

		expected := "Kaspi API error -1501: Device not found"
		if err.Error() != expected {
			t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("IsKaspiError identifies KaspiError", func(t *testing.T) {
		originalErr := &domain.KaspiError{
			StatusCode: -1501,
			Message:    "Device not found",
		}

		kaspiErr, ok := domain.IsKaspiError(originalErr)
		if !ok {
			t.Fatal("Expected IsKaspiError to return true for KaspiError")
		}

		if kaspiErr.StatusCode != -1501 {
			t.Errorf("Expected status code -1501, got %d", kaspiErr.StatusCode)
		}
	})

	t.Run("IsKaspiError unwraps wrapped KaspiError", func(t *testing.T) {
		originalErr := &domain.KaspiError{
			StatusCode: -1501,
			Message:    "Device not found",
		}
		wrappedErr := fmt.Errorf("operation failed: %w", originalErr)

		kaspiErr, ok := domain.IsKaspiError(wrappedErr)
		if !ok {
			t.Fatal("Expected IsKaspiError to return true for wrapped KaspiError")
		}

		if kaspiErr.StatusCode != -1501 {
			t.Errorf("Expected status code -1501, got %d", kaspiErr.StatusCode)
		}
	})

	t.Run("IsKaspiError returns false for non-KaspiError", func(t *testing.T) {
		regularErr := errors.New("regular error")

		_, ok := domain.IsKaspiError(regularErr)
		if ok {
			t.Fatal("Expected IsKaspiError to return false for regular error")
		}
	})

	t.Run("IsKaspiError returns false for nil error", func(t *testing.T) {
		_, ok := domain.IsKaspiError(nil)
		if ok {
			t.Fatal("Expected IsKaspiError to return false for nil error")
		}
	})

}
