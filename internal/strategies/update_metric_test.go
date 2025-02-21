package strategies

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
)

// MockSaver — мок для Saver
type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(key string, value *domain.Metrics) {
	m.Called(key, value)
}

// MockGetter — мок для Getter
type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) *domain.Metrics {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.Metrics)
}

// MockKeyEncoder — мок для KeyEncoder
type MockKeyEncoder struct {
	mock.Mock
}

func (m *MockKeyEncoder) Encode(id, mtype string) string {
	args := m.Called(id, mtype)
	return args.String(0)
}

func TestUpdateGaugeMetricStrategy_NewMetric(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	keyEncoder := new(MockKeyEncoder)

	strategy := NewUpdateGaugeMetricStrategy(saver, getter, keyEncoder)

	metric := &domain.Metrics{
		ID:    "cpu",
		MType: "gauge",
		Value: floatPtr(42.5),
	}

	encodedKey := "cpu:gauge"

	keyEncoder.On("Encode", "cpu", "gauge").Return(encodedKey)
	getter.On("Get", encodedKey).Return(nil)
	saver.On("Save", encodedKey, metric).Return()

	updatedMetric := strategy.Update(metric)

	assert.Equal(t, metric, updatedMetric)
	keyEncoder.AssertExpectations(t)
	getter.AssertExpectations(t)
	saver.AssertExpectations(t)
}

func TestUpdateGaugeMetricStrategy_UpdateExistingMetric(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	keyEncoder := new(MockKeyEncoder)

	strategy := NewUpdateGaugeMetricStrategy(saver, getter, keyEncoder)

	oldMetric := &domain.Metrics{
		ID:    "cpu",
		MType: "gauge",
		Value: floatPtr(30.0),
	}
	newMetric := &domain.Metrics{
		ID:    "cpu",
		MType: "gauge",
		Value: floatPtr(50.0),
	}

	encodedKey := "cpu:gauge"

	keyEncoder.On("Encode", "cpu", "gauge").Return(encodedKey)
	getter.On("Get", encodedKey).Return(oldMetric)
	saver.On("Save", encodedKey, mock.MatchedBy(func(m *domain.Metrics) bool {
		return *m.Value == *newMetric.Value
	})).Return()

	updatedMetric := strategy.Update(newMetric)

	assert.Equal(t, *newMetric.Value, *updatedMetric.Value)
	keyEncoder.AssertExpectations(t)
	getter.AssertExpectations(t)
	saver.AssertExpectations(t)
}

func TestUpdateCounterMetricStrategy_NewMetric(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	keyEncoder := new(MockKeyEncoder)

	strategy := NewUpdateCounterMetricStrategy(saver, getter, keyEncoder)

	metric := &domain.Metrics{
		ID:    "requests",
		MType: "counter",
		Delta: intPtr(10),
	}

	encodedKey := "requests:counter"

	keyEncoder.On("Encode", "requests", "counter").Return(encodedKey)
	getter.On("Get", encodedKey).Return(nil)
	saver.On("Save", encodedKey, metric).Return()

	updatedMetric := strategy.Update(metric)

	assert.Equal(t, metric, updatedMetric)
	keyEncoder.AssertExpectations(t)
	getter.AssertExpectations(t)
	saver.AssertExpectations(t)
}

func TestUpdateCounterMetricStrategy_UpdateExistingMetric(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	keyEncoder := new(MockKeyEncoder)

	strategy := NewUpdateCounterMetricStrategy(saver, getter, keyEncoder)

	oldMetric := &domain.Metrics{
		ID:    "requests",
		MType: "counter",
		Delta: intPtr(5),
	}
	newMetric := &domain.Metrics{
		ID:    "requests",
		MType: "counter",
		Delta: intPtr(10),
	}
	updatedDelta := *oldMetric.Delta + *newMetric.Delta

	encodedKey := "requests:counter"

	keyEncoder.On("Encode", "requests", "counter").Return(encodedKey)
	getter.On("Get", encodedKey).Return(oldMetric)
	saver.On("Save", encodedKey, mock.MatchedBy(func(m *domain.Metrics) bool {
		return *m.Delta == updatedDelta
	})).Return()

	updatedMetric := strategy.Update(newMetric)

	assert.Equal(t, updatedDelta, *updatedMetric.Delta)
	keyEncoder.AssertExpectations(t)
	getter.AssertExpectations(t)
	saver.AssertExpectations(t)
}

// Вспомогательные функции для указателей
func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int64) *int64 {
	return &i
}
