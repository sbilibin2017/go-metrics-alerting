package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/internal/validators"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking UpdateMetricsService
type MockUpdateMetricsService struct {
	mock.Mock
}

func (m *MockUpdateMetricsService) Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.UpdateMetricsResponse), args.Error(1)
}

// MockGetMetricValueService - Мок для сервиса получения метрик.
type MockGetMetricValueService struct {
	mock.Mock
}

// GetMetricValue - Имитация метода GetMetricValue.
func (m *MockGetMetricValueService) GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.GetMetricValueResponse), args.Error(1)
}

// MockGetAllMetricValuesService - мок для GetAllMetricValuesService
type MockGetAllMetricValuesService struct {
	mock.Mock
}

func (m *MockGetAllMetricValuesService) GetAllMetricValues() []*types.GetMetricValueResponse {
	args := m.Called()
	return args.Get(0).([]*types.GetMetricValueResponse)
}

// Test for invalid JSON body (parsing error)
func TestUpdateMetricsHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)
	invalidReqBody := "{invalid-json}"
	req, _ := http.NewRequest("POST", "/update-metrics", bytes.NewReader([]byte(invalidReqBody)))
	recorder := httptest.NewRecorder()
	r.POST("/update-metrics", UpdateMetricsHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "invalid character")
	assert.Contains(t, respBody["error"], "looking for beginning of object key string")
	mockService.AssertNotCalled(t, "Update", mock.Anything)
}

// Test for validation errors without service call
func TestUpdateMetricsHandler_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)
	tests := []struct {
		name           string
		reqBody        types.UpdateMetricsRequest
		expectedStatus int
		expectedResp   map[string]interface{}
		expectedError  error
	}{
		{
			name: "empty ID",
			reqBody: types.UpdateMetricsRequest{
				ID:    "",
				MType: "counter",
				Delta: new(int64),
				Value: new(float64),
			},
			expectedStatus: http.StatusNotFound,
			expectedResp: map[string]interface{}{
				"error": validators.ErrEmptyID.Error(),
			},
			expectedError: validators.ErrEmptyID,
		},
		{
			name: "invalid MType",
			reqBody: types.UpdateMetricsRequest{
				ID:    "1",
				MType: "invalid-type",
				Delta: new(int64),
				Value: new(float64),
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": validators.ErrInvalidMType.Error(),
			},
			expectedError: validators.ErrInvalidMType,
		},
		{
			name: "missing Delta for counter",
			reqBody: types.UpdateMetricsRequest{
				ID:    "1",
				MType: "counter",
				Delta: nil,
				Value: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": validators.ErrInvalidDelta.Error(),
			},
			expectedError: validators.ErrInvalidDelta,
		},
		{
			name: "missing Value for gauge",
			reqBody: types.UpdateMetricsRequest{
				ID:    "1",
				MType: "gauge",
				Delta: nil,
				Value: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": validators.ErrInvalidValue.Error(),
			},
			expectedError: validators.ErrInvalidValue,
		},
	}
	r.POST("/update-metrics", UpdateMetricsHandler(mockService))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/update-metrics", bytes.NewReader(reqBodyBytes))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			var respBody map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &respBody)
			assert.Equal(t, tt.expectedResp, respBody)
			mockService.AssertNotCalled(t, "Update", mock.Anything)
		})
	}
}

// Test for successful request and service call
func TestUpdateMetricsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)
	delta := int64(1)
	value := 2.5
	mockService.On("Update", mock.Anything).Return(&types.UpdateMetricsResponse{
		UpdateMetricsRequest: types.UpdateMetricsRequest{
			ID:    "1",
			MType: "counter",
			Delta: &delta,
			Value: &value,
		},
	}, nil)
	reqBody := types.UpdateMetricsRequest{
		ID:    "1",
		MType: "counter",
		Delta: &delta,
		Value: &value,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/update-metrics", bytes.NewReader(reqBodyBytes))
	recorder := httptest.NewRecorder()
	r.POST("/update-metrics", UpdateMetricsHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Equal(t, map[string]interface{}{
		"id":    "1",
		"mtype": "counter",
		"delta": float64(1),
		"value": 2.5,
	}, respBody)
	mockService.AssertCalled(t, "Update", mock.Anything)
}

// Test for service error during update
func TestUpdateMetricsHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)
	delta := int64(1)
	value := 2.5
	reqBody := types.UpdateMetricsRequest{
		ID:    "1",
		MType: "counter",
		Delta: &delta,
		Value: &value,
	}
	mockService.On("Update", mock.Anything).Return(nil, errors.New("internal server error"))
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/update-metrics", bytes.NewReader(reqBodyBytes))
	recorder := httptest.NewRecorder()
	r.POST("/update-metrics", UpdateMetricsHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "internal server error")
	mockService.AssertCalled(t, "Update", mock.Anything)
}

// Test for invalid JSON body (parsing error)
func TestGetMetricValueHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)
	invalidReqBody := "{invalid-json}"
	req, _ := http.NewRequest("POST", "/get-metric-value", bytes.NewReader([]byte(invalidReqBody)))
	recorder := httptest.NewRecorder()
	r.POST("/get-metric-value", GetMetricValueHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "invalid character")
	mockService.AssertNotCalled(t, "GetMetricValue", mock.Anything)
}

// Test for validation errors (ID and MType validation)
func TestGetMetricValueHandler_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)
	tests := []struct {
		name           string
		reqBody        types.GetMetricValueRequest
		expectedStatus int
		expectedResp   map[string]interface{}
		expectedError  error
	}{
		{
			name: "empty ID",
			reqBody: types.GetMetricValueRequest{
				ID:    "",
				MType: "counter",
			},
			expectedStatus: http.StatusNotFound,
			expectedResp: map[string]interface{}{
				"error": validators.ErrEmptyID.Error(),
			},
			expectedError: validators.ErrEmptyID,
		},
		{
			name: "invalid MType",
			reqBody: types.GetMetricValueRequest{
				ID:    "1",
				MType: "invalid-type",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": validators.ErrInvalidMType.Error(),
			},
			expectedError: validators.ErrInvalidMType,
		},
	}
	r.POST("/get-metric-value", GetMetricValueHandler(mockService))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/get-metric-value", bytes.NewReader(reqBodyBytes))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			var respBody map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &respBody)
			assert.Equal(t, tt.expectedResp, respBody)
			mockService.AssertNotCalled(t, "GetMetricValue", mock.Anything)
		})
	}
}

// Тест на ошибку сервиса
func TestGetMetricValueHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)
	reqBody := types.GetMetricValueRequest{
		ID:    "1",
		MType: "counter",
	}
	mockService.On("GetMetricValue", mock.Anything).Return(nil, errors.New("internal server error"))
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/get-metric-value", bytes.NewReader(reqBodyBytes))
	recorder := httptest.NewRecorder()
	r.POST("/get-metric-value", GetMetricValueHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "internal server error")
	mockService.AssertCalled(t, "GetMetricValue", mock.Anything)
}

// Тест на успешный ответ
func TestGetMetricValueHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)
	reqBody := types.GetMetricValueRequest{
		ID:    "1",
		MType: "counter",
	}
	mockResp := types.GetMetricValueResponse{
		ID:    "1",
		Value: "10",
	}
	mockService.On("GetMetricValue", mock.Anything).Return(&mockResp, nil)
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/value", bytes.NewReader(reqBodyBytes))
	recorder := httptest.NewRecorder()
	r.POST("/value", GetMetricValueHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	var respBody types.GetMetricValueResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "1", respBody.ID)
	assert.Equal(t, "10", respBody.Value)
	mockService.AssertCalled(t, "GetMetricValue", mock.Anything)
}

func TestGetAllMetricValuesHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetAllMetricValuesService)
	metrics := []*types.GetMetricValueResponse{
		{ID: "1", Value: "10"},
		{ID: "2", Value: "20"},
	}
	mockService.On("GetAllMetricValues").Return(metrics)
	req, _ := http.NewRequest("GET", "/metrics", nil)
	recorder := httptest.NewRecorder()
	r.GET("/metrics", GetAllMetricValuesHandler(mockService))
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<li>1: 10</li>")
	assert.Contains(t, recorder.Body.String(), "<li>2: 20</li>")
	mockService.AssertCalled(t, "GetAllMetricValues")
}
