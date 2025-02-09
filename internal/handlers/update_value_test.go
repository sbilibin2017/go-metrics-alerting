package handlers

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Мок-сервис, который реализует интерфейс UpdateValueService
type MockUpdateValueService struct{}

func (m *MockUpdateValueService) UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	// Логика мок-обработки
	if req.Name == "error" {
		return &apierror.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	if req.Type == "invalid" {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Unsupported metric type",
		}
	}
	return nil
}

func setupRouter(svc UpdateValueService) *gin.Engine {
	r := gin.Default()
	RegisterUpdateValueHandler(r, svc)
	return r
}

func TestUpdateValueHandler_Success(t *testing.T) {
	// Arrange
	mockService := &MockUpdateValueService{}
	r := setupRouter(mockService)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/counter/metric1/10", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Metric updated", w.Body.String())
}

func TestUpdateValueHandler_InvalidMetricType(t *testing.T) {
	// Arrange
	mockService := &MockUpdateValueService{}
	r := setupRouter(mockService)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/invalid/metric1/10", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Unsupported metric type")
}

func TestUpdateValueHandler_InternalServerError(t *testing.T) {
	// Arrange
	mockService := &MockUpdateValueService{}
	r := setupRouter(mockService)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/counter/error/10", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestUpdateValueHandler_RouteNotFound(t *testing.T) {
	// Arrange
	mockService := &MockUpdateValueService{}
	r := setupRouter(mockService)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/unknown/route", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Route not found", w.Body.String())
}

type MockUpdateValueServiceWithError struct{}

func (m *MockUpdateValueServiceWithError) UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	return errors.New("unexpected error") // Имитация неожиданной ошибки
}

func TestUpdateValueHandler_InternalServerError_GenericError(t *testing.T) {
	// Arrange
	mockService := &MockUpdateValueServiceWithError{}
	r := setupRouter(mockService)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/counter/metric1/10", nil)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}
