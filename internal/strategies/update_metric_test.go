package strategies

import (
	"testing"

	"go-metrics-alerting/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Моки для интерфейса Saver
type SaverMock struct {
	mock.Mock
}

func (s *SaverMock) Save(id string, value string) bool {
	args := s.Called(id, value)
	return args.Bool(0)
}

// Моки для интерфейса Getter
type GetterMock struct {
	mock.Mock
}

func (g *GetterMock) Get(id string) (string, bool) {
	args := g.Called(id)
	return args.String(0), args.Bool(1)
}

func TestUpdateGaugeStrategy_SaveFailure(t *testing.T) {
	// Создаем мок для Saver
	saverMock := new(SaverMock)
	metric := &domain.Metric{
		ID:    "metric1",
		Value: "100",
	}

	// Устанавливаем ожидания
	saverMock.On("Save", metric.ID, metric.Value).Return(false) // Ошибка сохранения

	// Создаем стратегию обновления
	strategy := &UpdateGaugeStrategy{saver: saverMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.False(t, ok, "Update should fail due to save failure")
	assert.Nil(t, updatedMetric, "Metric should be nil after failed update")
	saverMock.AssertExpectations(t)
}

func TestUpdateCounterStrategy_GetFailure(t *testing.T) {
	// Создаем моки для Saver и Getter
	saverMock := new(SaverMock)
	getterMock := new(GetterMock)

	metric := &domain.Metric{
		ID:    "metric2",
		Value: "50",
	}

	// Устанавливаем ожидания для Getter
	getterMock.On("Get", metric.ID).Return("", false) // Не удается получить значение

	// Создаем стратегию обновления
	strategy := &UpdateCounterStrategy{saver: saverMock, getter: getterMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.False(t, ok, "Update should fail due to getter failure")
	assert.Nil(t, updatedMetric, "Metric should be nil after failed update")
	getterMock.AssertExpectations(t)
}

func TestUpdateCounterStrategy_ParseCurrentValueFailure(t *testing.T) {
	// Создаем моки для Saver и Getter
	saverMock := new(SaverMock)
	getterMock := new(GetterMock)

	metric := &domain.Metric{
		ID:    "metric3",
		Value: "50",
	}

	// Устанавливаем ожидания для Getter
	getterMock.On("Get", metric.ID).Return("invalid", true) // Неверное значение для текущей метрики

	// Создаем стратегию обновления
	strategy := &UpdateCounterStrategy{saver: saverMock, getter: getterMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.False(t, ok, "Update should fail due to parse failure of current value")
	assert.Nil(t, updatedMetric, "Metric should be nil after failed update")
	getterMock.AssertExpectations(t)
}

func TestUpdateCounterStrategy_ParseNewValueFailure(t *testing.T) {
	// Создаем моки для Saver и Getter
	saverMock := new(SaverMock)
	getterMock := new(GetterMock)

	metric := &domain.Metric{
		ID:    "metric4",
		Value: "invalid", // Неверное значение для новой метрики
	}

	// Устанавливаем ожидания для Getter
	getterMock.On("Get", metric.ID).Return("100", true) // Текущее значение = 100

	// Создаем стратегию обновления
	strategy := &UpdateCounterStrategy{saver: saverMock, getter: getterMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.False(t, ok, "Update should fail due to parse failure of new value")
	assert.Nil(t, updatedMetric, "Metric should be nil after failed update")
	getterMock.AssertExpectations(t)
}

func TestUpdateCounterStrategy_SaveFailure(t *testing.T) {
	// Создаем моки для Saver и Getter
	saverMock := new(SaverMock)
	getterMock := new(GetterMock)

	metric := &domain.Metric{
		ID:    "metric5",
		Value: "50",
	}

	// Устанавливаем ожидания для Getter
	getterMock.On("Get", metric.ID).Return("100", true) // Текущее значение = 100

	// Устанавливаем ожидания для Saver
	saverMock.On("Save", metric.ID, "150").Return(false) // Ошибка сохранения

	// Создаем стратегию обновления
	strategy := &UpdateCounterStrategy{saver: saverMock, getter: getterMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.False(t, ok, "Update should fail due to save failure")
	assert.Nil(t, updatedMetric, "Metric should be nil after failed update")
	getterMock.AssertExpectations(t)
	saverMock.AssertExpectations(t)
}
func TestUpdateGaugeStrategy_Success(t *testing.T) {
	// Создаем моки для Saver
	saverMock := new(SaverMock)
	metric := &domain.Metric{
		ID:    "metric1",
		Value: "100",
	}

	// Устанавливаем ожидания
	saverMock.On("Save", metric.ID, metric.Value).Return(true) // Успешное сохранение

	// Создаем стратегию обновления
	strategy := &UpdateGaugeStrategy{saver: saverMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.True(t, ok, "Update should succeed")
	assert.NotNil(t, updatedMetric, "Updated metric should not be nil")
	assert.Equal(t, metric.ID, updatedMetric.ID, "Metric IDs should match")
	assert.Equal(t, metric.Value, updatedMetric.Value, "Metric values should match")
	saverMock.AssertExpectations(t)
}

func TestUpdateCounterStrategy_Success(t *testing.T) {
	// Создаем моки для Saver и Getter
	saverMock := new(SaverMock)
	getterMock := new(GetterMock)

	metric := &domain.Metric{
		ID:    "metric2",
		Value: "50",
	}

	// Устанавливаем ожидания для Getter
	getterMock.On("Get", metric.ID).Return("100", true) // Текущее значение = 100

	// Устанавливаем ожидания для Saver
	saverMock.On("Save", metric.ID, "150").Return(true) // Успешное сохранение

	// Создаем стратегию обновления
	strategy := &UpdateCounterStrategy{saver: saverMock, getter: getterMock}

	// Вызываем метод обновления
	updatedMetric, ok := strategy.Update(metric)

	// Проверяем результаты
	assert.True(t, ok, "Update should succeed")
	assert.NotNil(t, updatedMetric, "Updated metric should not be nil")
	assert.Equal(t, metric.ID, updatedMetric.ID, "Metric IDs should match")
	assert.Equal(t, "150", updatedMetric.Value, "Metric values should match after addition")
	getterMock.AssertExpectations(t)
	saverMock.AssertExpectations(t)
}
