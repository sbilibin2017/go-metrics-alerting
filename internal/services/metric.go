package services

import (
	"go-metrics-alerting/internal/types"
	"net/http"
)

// Setter интерфейс для записи данных.
type Saver interface {
	Save(key string, value string) error
}

// Getter интерфейс для чтения данных.
type Getter interface {
	Get(key string) (string, error)
}

// Ranger интерфейс для перебора данных.
type Ranger interface {
	Range(callback func(key string, value string) bool)
}

// MetricUpdateStrategy - интерфейс для стратегии обновления метрик.
type MetricUpdateStrategy interface {
	Update(req *types.MetricsRequest, currentValue string) (*types.MetricsRequest, error)
}

// IDValidator интерфейс для валидации ID.
type IDValidator interface {
	Validate(id string) bool
}

// MTypeValidator интерфейс для валидации типа метрики.
type MTypeValidator interface {
	Validate(mType string) bool
}

// DeltaValidator интерфейс для валидации Delta.
type DeltaValidator interface {
	Validate(mtype string, delta *int64) bool
}

// ValueValidator интерфейс для валидации Value.
type ValueValidator interface {
	Validate(mtype string, value *float64) bool
}

// UpdateMetricService - сервис для обновления метрик.
type UpdateMetricService struct {
	stringGetter   Getter
	stringSaver    Saver
	strategies     map[string]MetricUpdateStrategy
	idValidator    IDValidator
	mtypeValidator MTypeValidator
	deltaValidator DeltaValidator
	valueValidator ValueValidator
}

// NewUpdateMetricService создаёт новый сервис для обновления метрик.
func NewUpdateMetricService(
	stringGetter Getter,
	stringSaver Saver,
	strategies map[string]MetricUpdateStrategy,
	idValidator IDValidator,
	mtypeValidator MTypeValidator,
	deltaValidator DeltaValidator,
	valueValidator ValueValidator,
) *UpdateMetricService {
	return &UpdateMetricService{
		stringGetter:   stringGetter,
		stringSaver:    stringSaver,
		strategies:     strategies,
		idValidator:    idValidator,
		mtypeValidator: mtypeValidator,
		deltaValidator: deltaValidator,
		valueValidator: valueValidator,
	}
}

// UpdateMetric обновляет или сохраняет метрику в зависимости от её типа.
func (s *UpdateMetricService) UpdateMetric(req *types.MetricsRequest) (*types.MetricsResponse, *types.APIStatusResponse) {
	// Валидация ID
	if !s.idValidator.Validate(req.ID) {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusNotFound,
			Message: "Metric is not found",
		}
	}

	// Валидация MType
	if !s.mtypeValidator.Validate(req.MType) {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid metric type",
		}
	}

	// Валидация Delta для Counter метрик
	if !s.deltaValidator.Validate(req.MType, req.Delta) {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid delta for Counter metric",
		}
	}

	// Валидация Value для Gauge метрик
	if !s.valueValidator.Validate(req.MType, req.Value) {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid value for Gauge metric",
		}
	}

	// Получаем текущее значение метрики из хранилища
	currentValue, err := s.stringGetter.Get(req.ID)
	if err != nil {
		currentValue = "0"
	}

	// Выбираем стратегию для обновления в зависимости от типа метрики
	strategy, exists := s.strategies[string(req.MType)]
	if !exists {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusInternalServerError,
			Message: "No strategy found for metric type",
		}
	}

	// Используем стратегию для обновления метрики
	req, err = strategy.Update(req, currentValue)
	if err != nil {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusBadRequest,
			Message: "Metric is not updated",
		}
	}

	// Сохраняем обновленную метрику
	err = s.stringSaver.Save(req.ID, req.MType)
	if err != nil {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusInternalServerError,
			Message: "Metric is not saved",
		}
	}

	return &types.MetricsResponse{
		MetricsRequest: types.MetricsRequest{
			ID:    req.ID,
			MType: req.MType,
			Delta: req.Delta,
			Value: req.Value,
		},
	}, nil
}

// GetMetricService сервис для получения метрик.
type GetMetricService struct {
	stringGetter   Getter
	idValidator    IDValidator
	mtypeValidator MTypeValidator
}

// NewGetMetricService создаёт новый сервис для получения метрик.
func NewGetMetricService(
	stringGetter Getter,
	idValidator IDValidator,
	mtypeValidator MTypeValidator,
) *GetMetricService {
	return &GetMetricService{
		stringGetter:   stringGetter,
		idValidator:    idValidator,
		mtypeValidator: mtypeValidator,
	}
}

// GetMetric получает метрику по её ID.
func (s *GetMetricService) GetMetric(req *types.MetricValueRequest) (*string, *types.APIStatusResponse) {
	// Валидация ID
	if !s.idValidator.Validate(req.ID) {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusNotFound,
			Message: "Metric is not found",
		}
	}

	// Валидация MType
	if !s.mtypeValidator.Validate(req.MType) {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid metric type",
		}
	}

	// Получаем значение метрики
	value, err := s.stringGetter.Get(req.ID)
	if err != nil {
		return nil, &types.APIStatusResponse{
			Status:  http.StatusNotFound,
			Message: "Metric is not found",
		}
	}

	// Возвращаем успешный ответ
	return &value, nil
}

// ListMetricsService сервис для получения списка метрик.
type ListMetricsService struct {
	stringRanger Ranger
}

// NewListMetricsService создаёт новый сервис для получения списка метрик.
func NewListMetricsService(
	stringRanger Ranger,
) *ListMetricsService {
	return &ListMetricsService{
		stringRanger: stringRanger,
	}
}

// ListMetrics возвращает список всех метрик.
func (s *ListMetricsService) ListMetrics() []*types.MetricValueResponse {
	var metrics []*types.MetricValueResponse

	s.stringRanger.Range(func(key string, value string) bool {
		metrics = append(metrics, &types.MetricValueResponse{
			ID:    key,
			Value: value,
		})
		return true
	})

	return metrics
}
