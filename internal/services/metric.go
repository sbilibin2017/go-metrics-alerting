package services

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"strconv"
)

const (
	MetricEmptyString  string = ""
	DefaultMetricValue string = "0"
)

// MetricRepository — интерфейс хранилища метрик.
type MetricRepository interface {
	Save(metricType, metricName, metricValue string)
	Get(metricType, metricName string) (string, error)
	GetAll() [][]string
}

// MetricUpdateValueStrategy — интерфейс для стратегии обработки значения метрики.
type metricUpdateValueStrategy interface {
	update(currentValue, newValue string) string
}

// MetricService — сервис для работы с метриками.
type MetricService struct {
	repo MetricRepository
}

func NewMetricService(repo MetricRepository) *MetricService {
	return &MetricService{
		repo: repo,
	}
}

func (s *MetricService) UpdateMetric(req *types.UpdateMetricValueRequest) {
	currentValue, err := s.repo.Get(string(req.Type), req.Name)
	if err != nil || currentValue == MetricEmptyString {
		currentValue = DefaultMetricValue
	}
	strategy := updateStrategies[req.Type]
	value := strategy.update(currentValue, req.Value)
	s.repo.Save(string(req.Type), req.Name, value)
}

// GetMetricValue возвращает значение метрики по имени и типу.
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

// GetAllMetrics возвращает список всех метрик.
func (s *MetricService) ListMetrics() []*types.MetricResponse {
	metricsList := s.repo.GetAll()
	var metrics []*types.MetricResponse

	if len(metricsList) == 0 {
		return []*types.MetricResponse{}
	}

	for _, metric := range metricsList {
		metrics = append(metrics, &types.MetricResponse{
			Name:  metric[1],
			Value: metric[2],
		})
	}

	return metrics
}

var updateStrategies = map[types.MetricType]metricUpdateValueStrategy{
	types.Counter: &counterUpdateStrategy{},
	types.Gauge:   &gaugeUpdateStrategy{},
}

// CounterMetricStrategy — стратегия обработки значения для типа Counter.
type counterUpdateStrategy struct{}

func (s *counterUpdateStrategy) update(currentValue, newValue string) string {
	curVal, _ := strconv.ParseInt(currentValue, 10, 64) // Игнорируем ошибку, так как валидация на это уже была
	newVal, _ := strconv.ParseInt(newValue, 10, 64)     // Игнорируем ошибку, так как валидация на это уже была
	newVal += curVal
	return strconv.FormatInt(newVal, 10)
}

// GaugeMetricStrategy — стратегия обработки значения для типа Gauge.
type gaugeUpdateStrategy struct{}

func (s *gaugeUpdateStrategy) update(currentValue, newValue string) string {
	newVal, _ := strconv.ParseFloat(newValue, 64) // Игнорируем ошибку, так как валидация на это уже была
	return strconv.FormatFloat(newVal, 'f', -1, 64)
}
