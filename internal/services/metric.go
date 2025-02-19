package services

import (
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/types"
	"net/http"
)

// Validator interface defines the contract for all validators
type IDValidator interface {
	Validate(id string) error
}

// MTypeValidator interface for validating metric types
type MTypeValidator interface {
	Validate(mType types.MType) error
}

// DeltaValidator interface for validating delta values
type DeltaValidator interface {
	Validate(mType types.MType, delta *int64) error
}

// ValueValidator interface for validating value of gauges
type ValueValidator interface {
	Validate(mType types.MType, value *float64) error
}

// CounterValueValidator interface for validating counter values
type CounterValueValidator interface {
	Validate(value string) error
}

// GaugeValueValidator interface for validating gauge values
type GaugeValueValidator interface {
	Validate(value string) error
}

// Int64ParserFormatter interface defines methods for parsing and formatting int64
type Int64Formatter interface {
	Parse(value string) (int64, error)
	Format(value int64) string
}

// Float64ParserFormatter interface defines methods for parsing and formatting float64
type Float64Formatter interface {
	Parse(value string) (float64, error)
	Format(value float64) string
}

// Saver interface for the Save method
type Saver interface {
	Save(key, value string) bool
}

// Getter interface for the Get method
type Getter interface {
	Get(key string) (string, bool)
}

// Ranger interface for the Range method
type Ranger interface {
	Range(callback func(key, value string) bool)
}

// Обновленный UpdateMetricsService
type UpdateMetricsService struct {
	Saver            Saver
	Getter           Getter
	IDValidator      IDValidator
	MTypeValidator   MTypeValidator
	DeltaValidator   DeltaValidator
	ValueValidator   ValueValidator
	Int64Formatter   Int64Formatter
	Float64Formatter Float64Formatter
}

func NewUpdateMetricsService(
	Saver Saver,
	Getter Getter,
	IDValidator IDValidator,
	MTypeValidator MTypeValidator,
	DeltaValidator DeltaValidator,
	ValueValidator ValueValidator,
	Int64Formatter Int64Formatter,
	Float64Formatter Float64Formatter,
) *UpdateMetricsService {
	logger.Logger.Info("UpdateMetricsService initialized")
	return &UpdateMetricsService{
		Saver:            Saver,
		Getter:           Getter,
		IDValidator:      IDValidator,
		MTypeValidator:   MTypeValidator,
		DeltaValidator:   DeltaValidator,
		ValueValidator:   ValueValidator,
		Int64Formatter:   Int64Formatter,
		Float64Formatter: Float64Formatter,
	}
}

func (s *UpdateMetricsService) UpdateMetricValue(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, *types.APIErrorResponse) {
	// Валидация входных данных
	if err := s.IDValidator.Validate(req.ID); err != nil {
		return nil, &types.APIErrorResponse{
			Status:  http.StatusNotFound, // 404 Not Found - ID не найден
			Message: "Metric with the given ID not found",
		}
	}
	if err := s.MTypeValidator.Validate(req.MType); err != nil {
		return nil, &types.APIErrorResponse{
			Status:  http.StatusBadRequest, // 400 Bad Request
			Message: "Invalid metric type",
		}
	}
	if err := s.DeltaValidator.Validate(req.MType, req.Delta); err != nil {
		return nil, &types.APIErrorResponse{
			Status:  http.StatusBadRequest, // 400 Bad Request
			Message: "Invalid delta value",
		}
	}
	if err := s.ValueValidator.Validate(req.MType, req.Value); err != nil {
		return nil, &types.APIErrorResponse{
			Status:  http.StatusBadRequest, // 400 Bad Request
			Message: "Invalid value",
		}
	}

	switch req.MType {
	case types.Gauge:
		s.Saver.Save(req.ID, s.Float64Formatter.Format(*req.Value))
	case types.Counter:
		existingValue := int64(0)
		if existingValueStr, exists := s.Getter.Get(req.ID); exists {
			existingValue, _ = s.Int64Formatter.Parse(existingValueStr)
		}
		*req.Delta += existingValue
		s.Saver.Save(req.ID, s.Int64Formatter.Format(*req.Delta))
	}

	return &types.UpdateMetricsResponse{UpdateMetricsRequest: *req}, nil
}

func (s *UpdateMetricsService) ParseMetricValues(mtype, valueStr string) (*float64, *int64, *types.APIErrorResponse) {
	var value *float64
	var delta *int64

	// Преобразуем параметры в нужные типы
	if mtype == string(types.Counter) {
		parsedValue, err := s.Int64Formatter.Parse(valueStr)
		if err != nil {
			return nil, nil, &types.APIErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid value format for Counter",
			}
		}
		delta = &parsedValue
	} else if mtype == string(types.Gauge) {
		parsedValue, err := s.Float64Formatter.Parse(valueStr)
		if err != nil {
			return nil, nil, &types.APIErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid value format for Gauge",
			}
		}
		value = &parsedValue
	}
	return value, delta, nil
}

// Сервис получения метрики по ID
type GetMetricValueService struct {
	Getter         Getter
	IDValidator    IDValidator
	MTypeValidator MTypeValidator
}

func NewGetMetricValueService(getter Getter, emptyValidator IDValidator, mtypeValidator MTypeValidator) *GetMetricValueService {
	return &GetMetricValueService{
		Getter:         getter,
		IDValidator:    emptyValidator,
		MTypeValidator: mtypeValidator,
	}
}

func (s *GetMetricValueService) GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, *types.APIErrorResponse) {
	if err := s.IDValidator.Validate(req.ID); err != nil {
		return nil, &types.APIErrorResponse{
			Status:  http.StatusNotFound, // 404 Not Found - ID не найден
			Message: "Metric with the given ID not found",
		}
	}
	if err := s.MTypeValidator.Validate(req.MType); err != nil {
		return nil, &types.APIErrorResponse{
			Status:  http.StatusBadRequest, // 400 Bad Request
			Message: "Invalid metric type",
		}
	}

	if valueStr, exists := s.Getter.Get(req.ID); exists {
		return &types.GetMetricValueResponse{ID: req.ID, Value: valueStr}, nil
	}
	return nil, &types.APIErrorResponse{
		Status:  http.StatusNotFound, // 404 Not Found - метрика не найдена
		Message: "Metric not found",
	}
}

// Сервис получения всех метрик
type GetAllMetricValuesService struct {
	ranger Ranger
}

func NewGetAllMetricValuesService(ranger Ranger) *GetAllMetricValuesService {
	return &GetAllMetricValuesService{ranger: ranger}
}

func (s *GetAllMetricValuesService) GetAllMetricValues() ([]*types.GetMetricValueResponse, *types.APIErrorResponse) {
	var metrics []*types.GetMetricValueResponse
	s.ranger.Range(func(key, value string) bool {
		metrics = append(metrics, &types.GetMetricValueResponse{ID: key, Value: value})
		return true
	})
	return metrics, nil
}
