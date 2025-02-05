package apierror

import (
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name            string
		apiErr          *APIError
		expectedMessage string
	}{
		{
			name: "Not Found error",
			apiErr: &APIError{
				Code:    http.StatusNotFound,
				Message: "Not Found",
			},
			expectedMessage: "Not Found",
		},
		{
			name: "Internal Server Error",
			apiErr: &APIError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			expectedMessage: "Internal Server Error",
		},
		{
			name: "Bad Request error",
			apiErr: &APIError{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
			},
			expectedMessage: "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := tt.apiErr.Error()
			if message != tt.expectedMessage {
				t.Errorf("Expected message '%v', but got '%v'", tt.expectedMessage, message)
			}
		})
	}
}

func TestAPIError_ToResponse(t *testing.T) {
	tests := []struct {
		name            string
		apiErr          *APIError
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Not Found error",
			apiErr: &APIError{
				Code:    http.StatusNotFound,
				Message: "Not Found",
			},
			expectedCode:    http.StatusNotFound,
			expectedMessage: "Not Found",
		},
		{
			name: "Internal Server Error",
			apiErr: &APIError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "Internal Server Error",
		},
		{
			name: "Bad Request error",
			apiErr: &APIError{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, message := tt.apiErr.ToResponse()
			if code != tt.expectedCode {
				t.Errorf("Expected code %v, but got %v", tt.expectedCode, code)
			}
			if message != tt.expectedMessage {
				t.Errorf("Expected message '%v', but got '%v'", tt.expectedMessage, message)
			}
		})
	}
}
