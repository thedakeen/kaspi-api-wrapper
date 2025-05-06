package domain

import (
	"errors"
	"fmt"
)

type KaspiError struct {
	StatusCode int
	Message    string
}

// Error implements error interface
func (e *KaspiError) Error() string {
	return fmt.Sprintf("Kaspi API error %d: %s", e.StatusCode, e.Message)
}

// IsKaspiError checks if an error is kaspi error
func IsKaspiError(err error) (*KaspiError, bool) {
	if err == nil {
		return nil, false
	}

	var kaspiErr *KaspiError
	ok := errors.As(err, &kaspiErr)

	return kaspiErr, ok
}
