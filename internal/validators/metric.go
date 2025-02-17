package validators

import (
	"go-metrics-alerting/internal/types"
)

// IDValidatorImpl реализует интерфейс IDValidator для проверки ID.
type IDValidator struct{}

func (v *IDValidator) Validate(id string) bool {
	return id == types.EmptyString
}

// MTypeValidatorImpl реализует интерфейс MTypeValidator для проверки типа метрики.
type MTypeValidator struct{}

func (v *MTypeValidator) Validate(mType string) bool {
	return mType != string(types.Counter) && mType != string(types.Gauge)
}

// DeltaValidatorImpl реализует интерфейс DeltaValidator для проверки Delta.
type DeltaValidator struct{}

func (v *DeltaValidator) Validate(mtype string, delta *int64) bool {
	return mtype == string(types.Counter) && delta == nil
}

// ValueValidatorImpl реализует интерфейс ValueValidator для проверки Value.
type ValueValidator struct{}

func (v *ValueValidator) Validate(mtype string, value *float64) bool {
	return mtype == string(types.Gauge) && value == nil
}
