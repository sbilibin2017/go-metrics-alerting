package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Моки для валидаторов и сервиса
type MockGetValueService struct {
	mock.Mock
}

func (m *MockGetValueService) GetMetricValue(ctx context.Context, req *types.GetMetricValueRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

type MockMetricTypeValidator struct {
	mock.Mock
}

func (m *MockMetricTypeValidator) Validate(metricType string) error {
	args := m.Called(metricType)
	return args.Error(0)
}

type MockMetricNameValidator struct {
	mock.Mock
}

func (m *MockMetricNameValidator) Validate(metricName string) error {
	args := m.Called(metricName)
	return args.Error(0)
}

func TestGetMetricValueHandler_ServiceError(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "cpu"
	metricName := "usage"

	// Мокаем валидаторы
	mockTypeValidator.On("Validate", metricType).Return(nil).Once()
	mockNameValidator.On("Validate", metricName).Return(nil).Once()

	// Мокаем ошибку в сервисе (например, APIError)
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}).Return("", &apierror.APIError{Code: http.StatusInternalServerError, Message: "Service error occurred"}).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType+"/"+metricName, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата: ошибка от APIError
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Service error occurred", w.Body.String())

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueHandler_UnknownError(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "cpu"
	metricName := "usage"

	// Мокаем валидаторы
	mockTypeValidator.On("Validate", metricType).Return(nil).Once()
	mockNameValidator.On("Validate", metricName).Return(nil).Once()

	// Мокаем ошибку в сервисе, которая не является APIError (например, другая ошибка)
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}).Return("", assert.AnError).Once() // используем стандартную ошибку

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType+"/"+metricName, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата: ошибка 500 с сообщением "Internal Server Error"
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueHandler_Success(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "cpu"
	metricName := "usage"
	metricValue := "75"

	// Мокаем валидаторы
	mockTypeValidator.On("Validate", metricType).Return(nil).Once()
	mockNameValidator.On("Validate", metricName).Return(nil).Once()

	// Мокаем сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}).Return(metricValue, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType+"/"+metricName, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, metricValue, w.Body.String())

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueHandler_TypeValidationError(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "invalid"
	metricName := "usage"

	// Мокаем валидатор
	mockTypeValidator.On("Validate", metricType).Return(assert.AnError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType+"/"+metricName, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueHandler_NameValidationError(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "cpu"
	metricName := "invalid"

	// Мокаем валидаторы
	mockTypeValidator.On("Validate", metricType).Return(nil).Once()
	mockNameValidator.On("Validate", metricName).Return(assert.AnError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType+"/"+metricName, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_Success(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "cpu"
	metricValue := "1024"

	// Мокаем валидатор
	mockTypeValidator.On("Validate", metricType).Return(nil).Once()

	// Мокаем сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: metricType,
	}).Return(metricValue, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, metricValue, w.Body.String())

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_TypeValidationError(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "invalid"

	// Мокаем валидатор
	mockTypeValidator.On("Validate", metricType).Return(assert.AnError).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_ServiceError(t *testing.T) {
	mockService := new(MockGetValueService)
	mockTypeValidator := new(MockMetricTypeValidator)
	mockNameValidator := new(MockMetricNameValidator)

	handler := &GetValueHandler{
		service:             mockService,
		metricTypeValidator: mockTypeValidator,
		metricNameValidator: mockNameValidator,
	}

	r := gin.Default()
	RegisterGetMetricValueHandler(r, handler)

	metricType := "cpu"

	// Мокаем валидатор
	mockTypeValidator.On("Validate", metricType).Return(nil).Once()

	// Мокаем ошибку в сервисе
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: metricType,
	}).Return("", &apierror.APIError{Code: http.StatusInternalServerError, Message: "Service error"}).Once()

	req, _ := http.NewRequest(http.MethodGet, "/value/"+metricType, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Service error", w.Body.String())

	// Проверка вызовов моков
	mockTypeValidator.AssertExpectations(t)
	mockNameValidator.AssertExpectations(t)
	mockService.AssertExpectations(t)
}
