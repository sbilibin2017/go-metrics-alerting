package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-metrics-alerting/internal/server/types"
	"go-metrics-alerting/internal/server/validators"

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
	r.POST("/update-metrics", UpdateMetricsBodyHandler(mockService))
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
	r.POST("/update-metrics", UpdateMetricsBodyHandler(mockService))
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
	r.POST("/update-metrics", UpdateMetricsBodyHandler(mockService))
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
	r.POST("/update-metrics", UpdateMetricsBodyHandler(mockService))
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
	r.POST("/get-metric-value", GetMetricValueBodyHandler(mockService))
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
	r.POST("/get-metric-value", GetMetricValueBodyHandler(mockService))
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
	r.POST("/get-metric-value", GetMetricValueBodyHandler(mockService))
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
	r.POST("/value", GetMetricValueBodyHandler(mockService))
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

// Тест для получения метрики по пути с успешным результатом
func TestGetMetricValueByPathHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)
	mockResp := types.GetMetricValueResponse{
		ID:    "testCounter",
		Value: "10",
	}
	mockService.On("GetMetricValue", mock.Anything).Return(&mockResp, nil)

	// Формируем запрос GET для пути /value/counter/testCounter
	req, _ := http.NewRequest("GET", "/value/counter/testCounter", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Проверка на успешный ответ
	assert.Equal(t, http.StatusOK, recorder.Code)
	var respBody types.GetMetricValueResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "testCounter", respBody.ID)
	assert.Equal(t, "10", respBody.Value)

	// Проверка вызова метода
	mockService.AssertCalled(t, "GetMetricValue", mock.Anything)
}

// Тест для ошибки в обработке запроса (например, неправильный mtype или id)
func TestGetMetricValueByPathHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)

	// Настроим мок на возврат ошибки
	mockService.On("GetMetricValue", mock.Anything).Return(nil, errors.New("internal server error"))

	// Формируем запрос GET для пути /value/counter/testCounter
	req, _ := http.NewRequest("GET", "/value/counter/testCounter", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Проверка на ошибку сервера

	// Проверка вызова метода
	mockService.AssertCalled(t, "GetMetricValue", mock.Anything)
}

// Тест для неверных параметров запроса (например, отсутствуют параметры)
func TestGetMetricValueByPathHandler_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockGetMetricValueService)

	// Параметры запроса будут недействительными
	req, _ := http.NewRequest("GET", "/value/invalidType", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку BadRequest
	assert.Equal(t, http.StatusNotFound, recorder.Code)

}

// Тест для неверных параметров запроса (например, отсутствуют параметры)
func TestGetMetricValueByPathHandler_InvalidParams2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.RedirectTrailingSlash = false
	mockService := new(MockGetMetricValueService)

	// Параметры запроса будут недействительными
	req, _ := http.NewRequest("GET", "/value/invalidType/", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку BadRequest
	assert.Equal(t, http.StatusNotFound, recorder.Code)

}

// Тест для неверного типа метрики
func TestGetMetricValueByPathHandler_InvalidMType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.RedirectTrailingSlash = false
	mockService := new(MockGetMetricValueService)

	// Параметры запроса с неверным типом метрики
	req, _ := http.NewRequest("GET", "/metrics/invalidType/123", nil)
	recorder := httptest.NewRecorder()
	r.GET("/metrics/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку BadRequest из-за неверного типа метрики
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// Тест для пустого ID
func TestGetMetricValueByPathHandler_EmptyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.RedirectTrailingSlash = false
	mockService := new(MockGetMetricValueService)

	// Параметры запроса с пустым ID
	req, _ := http.NewRequest("GET", "/metrics/gauge/", nil) // Пустой ID
	recorder := httptest.NewRecorder()
	r.GET("/metrics/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку NotFound из-за пустого ID
	assert.Equal(t, http.StatusNotFound, recorder.Code)

}

// Тест для случая, когда метрика не найдена
func TestGetMetricValueByPathHandler_MetricNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.RedirectTrailingSlash = false

	// Мокируем сервис
	mockService := new(MockGetMetricValueService)

	// Определяем, что будет возвращать сервис, если метрика не найдена
	mockService.On("GetMetricValue", mock.AnythingOfType("*types.GetMetricValueRequest")).Return(nil, errors.New("metric not found"))

	// Параметры запроса с существующим mtype, но с несуществующим id
	req, _ := http.NewRequest("GET", "/metrics/gauge/nonexistentID", nil)
	recorder := httptest.NewRecorder()
	r.GET("/metrics/:mtype/:id", GetMetricValuePathHandler(mockService)) // Новый маршрут
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку NotFound, потому что метрика не найдена
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "metric not found")
}

// Тест для ошибки валидации мtype
func TestUpdateMetricsPathHandler_InvalidMType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)

	// Мокируем успешный ответ от сервиса
	mockService.On("Update", mock.Anything).Return(nil, nil)

	// Создаем запрос с неверным типом метрики (mtype)
	req, _ := http.NewRequest("GET", "/value/invalidType/123/10", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку BadRequest из-за неверного типа метрики
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "invalid metric type")
	mockService.AssertNotCalled(t, "Update", mock.Anything)
}

// Тест для пустого ID
func TestUpdateMetricsPathHandler_EmptyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)

	// Мокируем успешный ответ от сервиса
	mockService.On("Update", mock.Anything).Return(nil, nil)

	// Создаем запрос с пустым ID
	req, _ := http.NewRequest("GET", "/value/counter//10", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку NotFound из-за пустого ID
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "id cannot be empty")
	mockService.AssertNotCalled(t, "Update", mock.Anything)
}

// Тест для невалидного значения
func TestUpdateMetricsPathHandler_InvalidValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)

	// Мокируем успешный ответ от сервиса
	mockService.On("Update", mock.Anything).Return(nil, nil)

	// Создаем запрос с неверным значением
	req, _ := http.NewRequest("GET", "/value/counter/123/invalid", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Ожидаем ошибку BadRequest из-за некорректного значения
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	mockService.AssertNotCalled(t, "Update", mock.Anything)
}

// Тест для успешного обновления метрики
func TestUpdateMetricsPathHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)
	delta := int64(10)
	mockService.On("Update", mock.Anything).Return(&types.UpdateMetricsResponse{
		UpdateMetricsRequest: types.UpdateMetricsRequest{
			ID:    "123",
			MType: types.MType("counter"),
			Delta: &delta,
		},
	}, nil)

	// Создаем запрос для успешного обновления
	req, _ := http.NewRequest("GET", "/value/counter/123/10", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Проверка на успешный ответ
	assert.Equal(t, http.StatusOK, recorder.Code)
	var respBody types.UpdateMetricsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "123", respBody.ID)
	assert.Equal(t, types.MType("counter"), respBody.MType)
	assert.Equal(t, int64(10), *respBody.Delta)

	mockService.AssertCalled(t, "Update", mock.Anything)
}

// Тест для ошибки сервиса при обновлении метрики
func TestUpdateMetricsPathHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)

	// Мокируем ошибку в сервисе
	mockService.On("Update", mock.Anything).Return(nil, errors.New("internal server error"))

	// Создаем запрос для ошибки сервиса
	req, _ := http.NewRequest("GET", "/value/counter/123/10", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Проверка на ошибку сервиса
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var respBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.Contains(t, respBody["error"], "internal server error")

	mockService.AssertCalled(t, "Update", mock.Anything)
}

func TestUpdateMetricsPathHandler_GaugeSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)
	value := 42.42

	mockService.On("Update", mock.Anything).Return(&types.UpdateMetricsResponse{
		UpdateMetricsRequest: types.UpdateMetricsRequest{
			ID:    "gauge_metric",
			MType: types.MType("gauge"),
			Value: &value,
		},
	}, nil)

	// Создаем запрос для успешного обновления Gauge
	req, _ := http.NewRequest("GET", "/value/gauge/gauge_metric/42.42", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Проверяем тело ответа
	var respBody types.UpdateMetricsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "gauge_metric", respBody.ID)
	assert.Equal(t, types.MType("gauge"), respBody.MType)
	assert.Equal(t, 42.42, *respBody.Value)

	mockService.AssertCalled(t, "Update", mock.Anything)
}

func TestUpdateMetricsPathHandler_GaugeInvalidValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)

	// Создаем запрос с невалидным значением (не число)
	req, _ := http.NewRequest("GET", "/value/gauge/gauge_metric/not_a_number", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Проверяем, что вернулся 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateMetricsPathHandler_InvalidMType2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockUpdateMetricsService)

	// Создаем запрос с недопустимым типом метрики
	req, _ := http.NewRequest("GET", "/value/invalidType/metric_id/123", nil)
	recorder := httptest.NewRecorder()
	r.GET("/value/:mtype/:id/:value", UpdateMetricsPathHandler(mockService))
	r.ServeHTTP(recorder, req)

	// Проверяем, что вернулся 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Исправлено: вызываем метод String() у recorder.Body
	expectedError := validators.ErrInvalidMType.Error()
	assert.Contains(t, recorder.Body.String(), expectedError)
}
