package handlers_test

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the service to be used in the handler test
type MockGetMetricPathService struct {
	mock.Mock
}

func (m *MockGetMetricPathService) GetMetric(id string, mtype domain.MetricType) (*domain.Metrics, error) {
	args := m.Called(id, mtype)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Metrics), args.Error(1)
}

func TestGetMetricRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.GetMetricRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     &types.GetMetricRequest{ID: "metric1", MType: "counter"},
			wantErr: false,
		},
		{
			name:    "empty ID",
			req:     &types.GetMetricRequest{ID: "", MType: "counter"},
			wantErr: true,
		},
		{
			name:    "empty MType",
			req:     &types.GetMetricRequest{ID: "metric1", MType: ""},
			wantErr: true,
		},
		{
			name:    "invalid MType",
			req:     &types.GetMetricRequest{ID: "metric1", MType: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetMetricPathHandler_Success(t *testing.T) {
	mockService := new(MockGetMetricPathService)
	metric := &domain.Metrics{
		MType: domain.Counter,
		Delta: new(int64),
	}
	*metric.Delta = 10

	mockService.On("GetMetric", "metric1", domain.Counter).Return(metric, nil)

	r := chi.NewRouter()
	r.Get("/metrics/{id}/{type}", handlers.GetMetricPathHandler(mockService))

	req := httptest.NewRequest(http.MethodGet, "/metrics/metric1/counter", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "10", w.Body.String())
	mockService.AssertExpectations(t)
}

func TestGetMetricPathHandler_InvalidMetricType(t *testing.T) {
	mockService := new(MockGetMetricPathService)

	r := chi.NewRouter()
	r.Get("/metrics/{id}/{type}", handlers.GetMetricPathHandler(mockService))

	req := httptest.NewRequest(http.MethodGet, "/metrics/metric1/invalid", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMetricPathHandler_ServiceError(t *testing.T) {
	mockService := new(MockGetMetricPathService)

	mockService.On("GetMetric", "metric1", domain.Counter).Return(nil, assert.AnError)

	r := chi.NewRouter()
	r.Get("/metrics/{id}/{type}", handlers.GetMetricPathHandler(mockService))

	req := httptest.NewRequest(http.MethodGet, "/metrics/metric1/counter", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
