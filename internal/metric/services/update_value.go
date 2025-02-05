package services

import (
	"errors"
	"go-metrics-alerting/internal/metric/handlers"
	"go-metrics-alerting/pkg/apierror"

	"net/http"
	"strconv"
)

// Ошибки для сверки значений
var (
	ErrInvalidGaugeValue     = errors.New("invalid gauge value")
	ErrInvalidCounterValue   = errors.New("invalid counter value")
	ErrMetricMismatch        = errors.New("metric value mismatch during update")
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
)

const (
	updateValueEmptyString string = ""
	defaultCurrentValue    string = "0"
)

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

// Интерфейс для работы с хранилищем метрик
type MetricStorage interface {
	Save(metricType string, metricName string, value string) error
	Get(metricType string, metricName string) (string, error)
}

// Сервис для обновления и получения значений метрик
type UpdateMetricValueService struct {
	metricRepository MetricStorage
}

// Новый сервис для работы с метриками
func NewUpdateMetricValueService(metricRepository MetricStorage) *UpdateMetricValueService {
	return &UpdateMetricValueService{
		metricRepository: metricRepository,
	}
}

// UpdateMetricValue обновляет значение метрики
func (s *UpdateMetricValueService) UpdateMetricValue(req *handlers.UpdateMetricValueRequest) error {
	// Получаем текущее значение метрики из хранилища
	currentValue, err := s.metricRepository.Get(req.Type, req.Name)
	if err != nil {
		// Если метрики нет в хранилище, используем значение по умолчанию
		currentValue = defaultCurrentValue
	}

	// Проверяем и обновляем значение метрики
	updatedValue, err := updateMetric(req.Type, currentValue, req.Value)
	if err != nil {
		// Передаем ошибку как строку
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// Сохраняем обновленное значение в хранилище
	return s.metricRepository.Save(req.Type, req.Name, updatedValue)
}

// updateMetric обновляет значение метрики в зависимости от типа
func updateMetric(metricType, currentValue, newValue string) (string, error) {
	switch metricType {
	case string(GaugeType):
		return updateGauge(currentValue, newValue)
	case string(CounterType):
		return updateCounter(currentValue, newValue)
	default:
		// Возвращаем ошибку в виде строки
		return updateValueEmptyString, ErrUnsupportedMetricType
	}
}

// updateGauge обновляет значение gauge-метрики
func updateGauge(_, newValue string) (string, error) {
	newVal, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		// Возвращаем ошибку в виде строки
		return updateValueEmptyString, ErrInvalidGaugeValue
	}
	return strconv.FormatFloat(newVal, 'f', -1, 64), nil
}

// updateCounter обновляет значение counter-метрики
func updateCounter(currentValue, newValue string) (string, error) {
	// Проверка текущего значения
	current, err := strconv.ParseInt(currentValue, 10, 64)
	if err != nil {
		// Возвращаем ошибку в виде строки
		return updateValueEmptyString, ErrInvalidCounterValue
	}

	// Проверка нового значения
	newVal, err := strconv.ParseInt(newValue, 10, 64)
	if err != nil {
		// Возвращаем ошибку в виде строки
		return updateValueEmptyString, ErrInvalidCounterValue
	}

	// Проверка на сверку значений (например, для счетчика)
	if current+newVal < current {
		// Возвращаем ошибку в виде строки (переполнение)
		return updateValueEmptyString, ErrMetricMismatch
	}

	return strconv.FormatInt(current+newVal, 10), nil
}
