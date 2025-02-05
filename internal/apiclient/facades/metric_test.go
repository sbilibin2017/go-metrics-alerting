package facades

import (
	"errors"
	"go-metrics-alerting/internal/apiclient"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRequester struct {
	mock.Mock
}

func (m *MockPostRequester) Post(url string, headers map[string]string) (*apiclient.APIResponse, error) {
	args := m.Called(url, headers)
	if args.Get(0) != nil {
		return args.Get(0).(*apiclient.APIResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestMetricFacade_UpdateMetric_Success(t *testing.T) {
	mockClient := new(MockPostRequester)
	facade := NewMetricFacade(mockClient, "http://localhost") // Передаем с полным URL

	// Ожидаем, что будет передан тот же URL с полным путем и без модификации
	mockClient.On("Post", "http://localhost/update/gauge/cpu/99", mock.AnythingOfType("map[string]string")).
		Return(&apiclient.APIResponse{StatusCode: http.StatusOK, Body: "OK"}, nil)

	err := facade.UpdateMetric("gauge", "cpu", "99")
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestMetricFacade_UpdateMetric_RequestFailure(t *testing.T) {
	mockClient := new(MockPostRequester)
	facade := NewMetricFacade(mockClient, "http://localhost")

	mockClient.On("Post", "http://localhost/update/gauge/cpu/99", mock.AnythingOfType("map[string]string")).
		Return(nil, errors.New("request failed"))

	err := facade.UpdateMetric("gauge", "cpu", "99")
	assert.ErrorIs(t, err, ErrRequestFailed)
	mockClient.AssertExpectations(t)
}

func TestMetricFacade_UpdateMetric_InvalidStatusCode(t *testing.T) {
	mockClient := new(MockPostRequester)
	facade := NewMetricFacade(mockClient, "http://localhost")

	mockClient.On("Post", "http://localhost/update/gauge/cpu/99", mock.AnythingOfType("map[string]string")).
		Return(&apiclient.APIResponse{StatusCode: http.StatusBadRequest, Body: "Bad Request"}, nil)

	err := facade.UpdateMetric("gauge", "cpu", "99")
	assert.ErrorIs(t, err, ErrInvalidStatusCode)
	mockClient.AssertExpectations(t)
}
