package services

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"strconv"
)

const (
	MetricEmptyString  = ""
	DefaultMetricValue = "0"
)

// MetricRepository — интерфейс хранилища метрик.
type MetricRepository interface {
	Save(metricType, metricName, metricValue string)
	Get(metricType, metricName string) (string, error)
	GetAll() [][]string
}

// MetricUpdateValueStrategy — интерфейс для стратегии обработки значения метрики.
type MetricUpdateValueStrategy interface {
	Update(currentValue, newValue string) string
}

// MetricService — сервис для работы с метриками.
type MetricService struct {
	repo       MetricRepository
	strategies map[types.MetricType]MetricUpdateValueStrategy
}

func NewMetricService(repo MetricRepository) *MetricService {
	return &MetricService{
		repo: repo,
		strategies: map[types.MetricType]MetricUpdateValueStrategy{
			types.Counter: &CounterMetricStrategy{},
			types.Gauge:   &GaugeMetricStrategy{},
		},
	}
}

func (s *MetricService) UpdateMetric(req *types.UpdateMetricValueRequest) {
	currentValue, err := s.repo.Get(string(req.Type), req.Name)
	if err != nil || currentValue == MetricEmptyString {
		currentValue = DefaultMetricValue
	}
	if strategy, ok := s.strategies[req.Type]; ok {
		value := strategy.Update(currentValue, req.Value)
		s.repo.Save(string(req.Type), req.Name, value)
	}
}

// GetMetric возвращает значение метрики по имени и типу.
func (s *MetricService) GetMetric(req *types.GetMetricValueRequest) (string, *types.APIErrorResponse) {
	currentValue, err := s.repo.Get(string(req.Type), req.Name)
	if err != nil {
		return MetricEmptyString, &types.APIErrorResponse{
			Code:    http.StatusNotFound,
			Message: "value not found",
		}
	}
	return currentValue, nil
}

// ListMetrics возвращает список всех метрик.
func (s *MetricService) ListMetrics() []*types.MetricResponse {
	metricsList := s.repo.GetAll()
	var metrics []*types.MetricResponse

	for _, metric := range metricsList {
		metrics = append(metrics, &types.MetricResponse{
			Name:  metric[1],
			Value: metric[2],
		})
	}

	return metrics
}

// CounterMetricStrategy — стратегия обработки значения для типа Counter.
type CounterMetricStrategy struct{}

func (s *CounterMetricStrategy) Update(currentValue, newValue string) string {
	curVal, _ := strconv.ParseInt(currentValue, 10, 64) // Игнорируем ошибку, так как валидация уже была
	newVal, _ := strconv.ParseInt(newValue, 10, 64)     // Игнорируем ошибку, так как валидация уже была
	newVal += curVal
	return strconv.FormatInt(newVal, 10)
}

// GaugeMetricStrategy — стратегия обработки значения для типа Gauge.
type GaugeMetricStrategy struct{}

func (s *GaugeMetricStrategy) Update(currentValue, newValue string) string {
	newVal, _ := strconv.ParseFloat(newValue, 64) // Игнорируем ошибку, так как валидация уже была
	return strconv.FormatFloat(newVal, 'f', -1, 64)
}
