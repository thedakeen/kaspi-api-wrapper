package http_test

import (
	"errors"
	httphandler "kaspi-api-wrapper/internal/api/http"
	"kaspi-api-wrapper/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleKaspiError(t *testing.T) {
	log := setupTestLogger()

	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Device not found",
			err:            &domain.KaspiError{StatusCode: -1501, Message: "Device not found"},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Device not found",
		},
		{
			name:           "Device not active",
			err:            &domain.KaspiError{StatusCode: -1502, Message: "Device not active"},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Device is not active",
		},
		{
			name:           "Device already registered",
			err:            &domain.KaspiError{StatusCode: -1503, Message: "Device already registered"},
			expectedStatus: http.StatusConflict,
			expectedMsg:    "Device is already registered to another trade point",
		},
		{
			name:           "Payment not found",
			err:            &domain.KaspiError{StatusCode: -1601, Message: "Payment not found"},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Payment not found",
		},
		{
			name:           "No trade points",
			err:            &domain.KaspiError{StatusCode: -14000002, Message: "No trade points"},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "No trade points available. Please create a trade point in the Kaspi Pay application",
		},
		{
			name:           "Service unavailable",
			err:            &domain.KaspiError{StatusCode: -999, Message: "Service unavailable"},
			expectedStatus: http.StatusServiceUnavailable,
			expectedMsg:    "Kaspi Pay service is temporarily unavailable",
		},
		{
			name:           "Unknown error",
			err:            &domain.KaspiError{StatusCode: -12345, Message: "Unknown error"},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unexpected error from payment system: Unknown error",
		},
		{
			name:           "Non-Kaspi error",
			err:            errors.New("Some other error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			httphandler.HandleError(recorder, tc.err, log)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d got %d", tc.expectedStatus, recorder.Code)
			}

			var resp httphandler.Response
			err := parseResponse(recorder, &resp)
			if err != nil {
				t.Fatalf("Failed to parce response: %v", err)
			}

			if resp.Success {
				t.Errorf("Expected success to be false, got true")
			}

			if resp.Error != tc.expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedMsg, resp.Error)
			}
		})
	}
}

func TestResponseHelpers(t *testing.T) {
	testCases := []struct {
		name           string
		helperFunc     func(http.ResponseWriter, string)
		expectedStatus int
	}{
		{
			name:           "BadRequestError",
			helperFunc:     httphandler.BadRequestError,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InternalServerError",
			helperFunc:     httphandler.InternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "NotFoundError",
			helperFunc:     httphandler.NotFoundError,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ConflictError",
			helperFunc:     httphandler.ConflictError,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "ForbiddenError",
			helperFunc:     httphandler.ForbiddenError,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ServiceUnavailableError",
			helperFunc:     httphandler.ServiceUnavailableError,
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "UnauthorizedError",
			helperFunc:     httphandler.UnauthorizedError,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			tc.helperFunc(recorder, "Test error message")

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, recorder.Code)
			}

			var resp httphandler.Response
			err := parseResponse(recorder, &resp)
			if err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if resp.Success {
				t.Errorf("Expected success to be false, got true")
			}

			if resp.Error != "Test error message" {
				t.Errorf("Expected error message 'Test error message', got '%s'", resp.Error)
			}
		})
	}
}

func TestDecodeJSONRequest(t *testing.T) {
	type TestRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("successfully decodes valid JSON", func(t *testing.T) {
		req, err := createRequest(http.MethodPost, "/test", TestRequest{
			Name: "John",
			Age:  30,
		})
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		var decoded TestRequest
		result := httphandler.DecodeJSONRequest(recorder, req, &decoded)

		if !result {
			t.Error("Expected DecodeJSONRequest to return true, got false")
		}

		if decoded.Name != "John" {
			t.Errorf("Expected Name John, got %s", decoded.Name)
		}

		if decoded.Age != 30 {
			t.Errorf("Expected Age 30, got %d", decoded.Age)
		}
	})

	t.Run("rejects invalid content type", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "text/plain")

		recorder := httptest.NewRecorder()

		var decoded TestRequest
		result := httphandler.DecodeJSONRequest(recorder, req, &decoded)

		if result {
			t.Error("Expected DecodeJSONRequest to return false, got true")
		}

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		expectedError := "Content-Type must be application/json"
		if resp.Error != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, resp.Error)
		}
	})

	t.Run("rejects malformed JSON", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/test",
			strings.NewReader(`{"name": "John", "age": "thirty"}`)) // age should be a number
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		var decoded TestRequest
		result := httphandler.DecodeJSONRequest(recorder, req, &decoded)

		if result {
			t.Error("Expected DecodeJSONRequest to return false, got true")
		}

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var resp httphandler.Response
		err = parseResponse(recorder, &resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Success {
			t.Errorf("Expected success to be false, got true")
		}

		if resp.Error == "" {
			t.Error("Expected non-empty error message")
		}
	})
}
