package routers

import (
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем зависимости
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Создаем роутер
	r := NewMetricRouter(metricService)

	// Проверяем, что r — это *gin.Engine
	assert.IsType(t, &gin.Engine{}, r, "NewMetricRouter должен возвращать *gin.Engine")
}

// Тест успешного получения метрики
func TestGetMetricHandler_Success(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.GET("/value/:type/:name", getMetricHandler(metricService))

	// Сначала добавим метрику, чтобы потом ее получить
	req := &types.UpdateMetricValueRequest{
		Type:  types.MetricType("gauge"),
		Name:  "metric1",
		Value: "100",
	}

	// Обновление метрики
	metricService.UpdateMetric(req)

	// Тест запроса на получение метрики
	request, _ := http.NewRequest("GET", "/value/gauge/metric1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	// Проверка успешного ответа
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "100")
}

// Тест на ошибку, если метрика не найдена
func TestGetMetricHandler_MetricNotFound(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.GET("/value/:type/:name", getMetricHandler(metricService))

	// Тест запроса на получение несуществующей метрики
	request, _ := http.NewRequest("GET", "/value/gauge/metric1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	// Проверка на ошибку 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Тест на успешное обновление метрики
func TestUpdateMetricHandler_Success(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.POST("/update/:type/:name/:value", updateMetricHandler(metricService))

	// Тест успешного обновления метрики
	request, _ := http.NewRequest("POST", "/update/gauge/metric1/200", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	// Проверка на успешный ответ
	assert.Equal(t, http.StatusOK, w.Code)
}

// Тест на получение списка метрик (пустой список)
func TestListMetricsHandler_Empty(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.GET("/", listMetricsHandler(metricService))

	// Тест запроса на получение списка метрик
	request, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	// Проверка на успешный ответ
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No metrics available")
}

func TestUpdateMetricHandler_ValidRequest(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.POST("/update/:type/:name/:value", updateMetricHandler(metricService))

	// Создаем правильный запрос на обновление метрики
	req := httptest.NewRequest("POST", "/update/gauge/metric1/100.5", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Проверяем, что ответ был успешным
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateMetricHandler_EmptyName(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.POST("/update/:type/:name/:value", updateMetricHandler(metricService))

	// Создаем запрос с пустым именем метрики
	req := httptest.NewRequest("POST", "/update/gauge//100.5", nil) // Пустое имя
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Проверяем, что произошла ошибка
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "metric name is required")
}

func TestUpdateMetricHandler_InvalidType(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.PUT("/update/:type/:name/:value", updateMetricHandler(metricService))

	// Создаем запрос с неверным типом метрики
	req := httptest.NewRequest("PUT", "/update/unknown/metric1/100.5", nil) // Неверный тип
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Проверяем, что произошла ошибка
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid metric type")
}

func TestUpdateMetricHandler_InvalidValueForGauge(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.PUT("/update/:type/:name/:value", updateMetricHandler(metricService))

	// Создаем запрос с неверным значением для метрики типа Gauge
	req := httptest.NewRequest("PUT", "/update/gauge/metric1/invalidValue", nil) // Неверное значение
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Проверяем, что произошла ошибка
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid gauge value")
}

func TestUpdateMetricHandler_InvalidValueForCounter(t *testing.T) {
	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Инициализация всех зависимостей
	memStorage := storage.NewMemStorage()
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Регистрируем роутер
	r := gin.Default()
	r.PUT("/update/:type/:name/:value", updateMetricHandler(metricService))

	// Создаем запрос с неверным значением для метрики типа Counter
	req := httptest.NewRequest("PUT", "/update/counter/metric1/invalidValue", nil) // Неверное значение
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Проверяем, что произошла ошибка
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid counter value")
}
