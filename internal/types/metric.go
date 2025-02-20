package types

import "go-metrics-alerting/internal/validators"

// UpdateMetricBodyRequest структура для обновления метрики.
type UpdateMetricBodyRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// Validate проверяет поля структуры UpdateMetricBodyRequest.
func (r *UpdateMetricBodyRequest) Validate() error {
	// Валидируем ID
	if err := validators.ValidateEmptyString(r.ID); err != nil {
		return err
	}

	// Валидируем тип метрики (MType)
	if err := validators.ValidateMType(r.MType); err != nil {
		return err
	}

	// Валидация Delta для типа MType "counter"
	if err := validators.ValidateDelta(r.MType, r.Delta); err != nil {
		return err
	}

	// Валидация Value для типа MType "gauge"
	if err := validators.ValidateValue(r.MType, r.Value); err != nil {
		return err
	}

	return nil
}

// UpdateMetricQueryRequest структура для запроса метрики.
type UpdateMetricQueryRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Validate проверяет поля структуры UpdateMetricQueryRequest.
func (r *UpdateMetricQueryRequest) Validate() error {
	// Валидируем Name
	if err := validators.ValidateEmptyString(r.Name); err != nil {
		return err
	}

	// Валидируем Type
	if err := validators.ValidateEmptyString(r.Type); err != nil {
		return err
	}

	// Валидируем Value
	if err := validators.ValidateEmptyString(r.Value); err != nil {
		return err
	}

	return nil
}

// GetMetricRequest структура для запроса метрик
type GetMetricBodyRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

// Validate проверяет поля структуры GetMetricRequest
func (r *GetMetricBodyRequest) Validate() error {
	// Валидируем поле ID
	if err := validators.ValidateEmptyString(r.ID); err != nil {
		return err
	}

	// Валидируем тип метрики (MType)
	if err := validators.ValidateMType(r.MType); err != nil {
		return err
	}

	return nil
}

// GetMetricRequest структура для ответа с метриками
type GetMetricBodyResponse struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Value string `json:"value"`
}
