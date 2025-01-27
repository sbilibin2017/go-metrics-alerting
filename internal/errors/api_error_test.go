package errors

import "testing"

func TestApiError_GetStatus(t *testing.T) {
	tests := []struct {
		name     string
		apiError *ApiError
		expected int
	}{
		{
			name:     "Status 200 OK",
			apiError: &ApiError{StatusCode: 200, Message: "OK"},
			expected: 200,
		},
		{
			name:     "Status 404 Not Found",
			apiError: &ApiError{StatusCode: 404, Message: "Not Found"},
			expected: 404,
		},
		{
			name:     "Status 500 Internal Server Error",
			apiError: &ApiError{StatusCode: 500, Message: "Internal Server Error"},
			expected: 500,
		},
		{
			name:     "Status 0 Invalid Status",
			apiError: &ApiError{StatusCode: 0, Message: "Invalid Status"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiError.Status() // Исправлено с Error() на Status()
			if got != tt.expected {
				t.Errorf("ApiError.Status() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestApiError_GetMessage(t *testing.T) {
	tests := []struct {
		name     string
		apiError *ApiError
		expected string
	}{
		{
			name:     "Message 'OK' for status 200",
			apiError: &ApiError{StatusCode: 200, Message: "OK"},
			expected: "OK",
		},
		{
			name:     "Message 'Not Found' for status 404",
			apiError: &ApiError{StatusCode: 404, Message: "Not Found"},
			expected: "Not Found",
		},
		{
			name:     "Message 'Internal Server Error' for status 500",
			apiError: &ApiError{StatusCode: 500, Message: "Internal Server Error"},
			expected: "Internal Server Error",
		},
		{
			name:     "Empty message",
			apiError: &ApiError{StatusCode: 404, Message: ""},
			expected: "",
		},
		{
			name:     "Very long message",
			apiError: &ApiError{StatusCode: 400, Message: "This is a very long error message that exceeds normal lengths."},
			expected: "This is a very long error message that exceeds normal lengths.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiError.Error()
			if got != tt.expected {
				t.Errorf("ApiError.Error() = %v, want %v", got, tt.expected) // Исправлено название метода в сообщении ошибки
			}
		})
	}
}
