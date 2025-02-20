package services

import (
	"testing"

	"go-metrics-alerting/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Моки для стратегий обновления
type UpdateMetricStrategyMock struct {
	mock.Mock
}

func (m *UpdateMetricStrategyMock) Update(metric *domain.Metric) (*domain.Metric, bool) {
	args := m.Called(metric)
	// Проверяем, что мы не возвращаем nil в случае ошибки
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*domain.Metric), args.Bool(1)
}

func TestUpdateMetricsService_UpdateMetricValue_Success(t *testing.T) {
	// Создаем тестовую метрику
	metric := &domain.Metric{
		ID:    "metric1",
		Value: "100",
	}

	// Мокаем стратегию
	strategyMock := new(UpdateMetricStrategyMock)
	strategyMock.On("Update", metric).Return(metric, true) // Стратегия успешно обновляет метрику

	// Создаем сервис с мокой стратегии
	service := &UpdateMetricsService{strategy: strategyMock}

	// Вызываем метод обновления
	updatedMetric, err := service.UpdateMetricValue(metric)

	// Проверяем, что ошибки нет и метрика обновлена корректно
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, metric.ID, updatedMetric.ID, "Metric IDs should match")
	assert.Equal(t, metric.Value, updatedMetric.Value, "Metric values should match")

	// Проверяем, что стратегия была вызвана
	strategyMock.AssertExpectations(t)
}

func TestUpdateMetricsService_UpdateMetricValue_Fail(t *testing.T) {
	// Создаем тестовую метрику
	metric := &domain.Metric{
		ID:    "metric1",
		Value: "100",
	}

	// Мокаем стратегию
	strategyMock := new(UpdateMetricStrategyMock)
	strategyMock.On("Update", metric).Return(nil, false) // Стратегия не обновляет метрику

	// Создаем сервис с мокой стратегии
	service := &UpdateMetricsService{strategy: strategyMock}

	// Вызываем метод обновления
	updatedMetric, err := service.UpdateMetricValue(metric)

	// Проверяем, что ошибка произошла и метрика не обновлена
	assert.Error(t, err, "Expected error")
	assert.Equal(t, ErrUpdateFailed, err, "Expected ErrUpdateFailed error")
	assert.Nil(t, updatedMetric, "Expected nil updated metric")

	// Проверяем, что стратегия была вызвана
	strategyMock.AssertExpectations(t)
}

// Мок для интерфейса Getter
type GetterMock struct {
	mock.Mock
}

func (m *GetterMock) Get(id string) (string, bool) {
	args := m.Called(id)
	return args.String(0), args.Bool(1)
}

// Мок для интерфейса KeyEncoder
type KeyEncoderMock struct {
	mock.Mock
}

func (m *KeyEncoderMock) Encode(id string, mtype string) string {
	args := m.Called(id, mtype)
	return args.String(0)
}

func TestGetMetricValue_Success(t *testing.T) {
	// Мокаем методы интерфейсов
	getterMock := new(GetterMock)
	encoderMock := new(KeyEncoderMock)

	// Определяем поведение моков
	getterMock.On("Get", "encoded_metric1_type1").Return("100", true) // Возвращаем значение метрики
	encoderMock.On("Encode", "metric1", "type1").Return("encoded_metric1_type1")

	// Создаем сервис с мокой
	service := &GetMetricValueService{
		getter:  getterMock,
		encoder: encoderMock,
	}

	// Вызываем метод
	result, err := service.GetMetricValue("metric1", domain.MType("type1"))

	// Проверяем, что ошибки нет и метрика получена корректно
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "metric1", result.ID, "Metric IDs should match")
	assert.Equal(t, "100", result.Value, "Metric values should match")
	assert.Equal(t, domain.MType("type1"), result.MType, "Metric types should match")

	// Проверяем, что методы Get и Encode были вызваны
	getterMock.AssertExpectations(t)
	encoderMock.AssertExpectations(t)
}

func TestGetMetricValue_Fail(t *testing.T) {
	// Мокаем методы интерфейсов
	getterMock := new(GetterMock)
	encoderMock := new(KeyEncoderMock)

	// Определяем поведение моков
	getterMock.On("Get", "encoded_metric1_type1").Return("", false) // Возвращаем false (метрика не найдена)
	encoderMock.On("Encode", "metric1", "type1").Return("encoded_metric1_type1")

	// Создаем сервис с мокой
	service := &GetMetricValueService{
		getter:  getterMock,
		encoder: encoderMock,
	}

	// Вызываем метод
	result, err := service.GetMetricValue("metric1", domain.MType("type1"))

	// Проверяем, что ошибка произошла
	assert.Error(t, err, "Expected error")
	assert.Nil(t, result, "Expected nil result")

	// Проверяем, что методы Get и Encode были вызваны
	getterMock.AssertExpectations(t)
	encoderMock.AssertExpectations(t)
}

// Мок для интерфейса Ranger
type RangerMock struct {
	mock.Mock
}

func (m *RangerMock) Range(callback func(id, value string) bool) {
	m.Called(callback)
	// Здесь можно определять, как метод Range будет обрабатывать вызов callback
	callback("encoded_metric1_type1", "100")
	callback("encoded_metric2_type2", "200")
	// Если мы случайно добавим больше метрик, убедимся, что callback не вызывается лишний раз.
}

// Мок для интерфейса KeyDecoder
type KeyDecoderMock struct {
	mock.Mock
}

func (m *KeyDecoderMock) Decode(key string) (id string, mtype string, ok bool) {
	args := m.Called(key)
	return args.String(0), args.String(1), args.Bool(2)
}

func TestGetAllMetricValues_Success(t *testing.T) {
	// Мокаем методы интерфейсов
	rangerMock := new(RangerMock)
	decoderMock := new(KeyDecoderMock)

	// Определяем поведение моков
	decoderMock.On("Decode", "encoded_metric1_type1").Return("metric1", "type1", true)
	decoderMock.On("Decode", "encoded_metric2_type2").Return("metric2", "type2", true)

	// Мокаем вызов метода Range
	rangerMock.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		// Вызываем callback с тестовыми значениями
		callback := args.Get(0).(func(id, value string) bool)
		callback("encoded_metric1_type1", "100")
		callback("encoded_metric2_type2", "200")
	}).Once()

	// Создаем сервис с мокой
	service := &GetAllMetricValuesService{
		ranger:  rangerMock,
		decoder: decoderMock,
	}

	// Вызываем метод
	metrics := service.GetAllMetricValues()

	assert.Equal(t, "metric1", metrics[0].ID, "Metric 1 ID should match")
	assert.Equal(t, "100", metrics[0].Value, "Metric 1 value should match")
	assert.Equal(t, domain.MType("type1"), metrics[0].MType, "Metric 1 type should match")

	assert.Equal(t, "metric2", metrics[1].ID, "Metric 2 ID should match")
	assert.Equal(t, "200", metrics[1].Value, "Metric 2 value should match")
	assert.Equal(t, domain.MType("type2"), metrics[1].MType, "Metric 2 type should match")

	// Проверяем, что методы Range и Decode были вызваны
	rangerMock.AssertExpectations(t)
	decoderMock.AssertExpectations(t)
}

func TestGetAllMetricValues_DecodeFail(t *testing.T) {
	// Мокаем методы интерфейсов
	rangerMock := new(RangerMock)
	decoderMock := new(KeyDecoderMock)

	// Мокаем вызов метода Range
	rangerMock.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(0).(func(id, value string) bool)
		// Передаем метрики, одна из которых не будет декодирована
		callback("encoded_metric1_type1", "100")
		callback("encoded_metric2_type2", "200")
	}).Once()

	// Мокаем поведение декодера, чтобы вторая метрика не была декодирована
	decoderMock.On("Decode", "encoded_metric1_type1").Return("metric1", "type1", true)
	decoderMock.On("Decode", "encoded_metric2_type2").Return("", "", false) // Ошибка декодирования

	// Создаем сервис с мокой
	service := &GetAllMetricValuesService{
		ranger:  rangerMock,
		decoder: decoderMock,
	}

	// Вызываем метод
	metrics := service.GetAllMetricValues()

	assert.Equal(t, "metric1", metrics[0].ID, "Metric 1 ID should match")
	assert.Equal(t, "100", metrics[0].Value, "Metric 1 value should match")
	assert.Equal(t, domain.MType("type1"), metrics[0].MType, "Metric 1 type should match")

	// Проверяем, что методы Range и Decode были вызваны
	rangerMock.AssertExpectations(t)
	decoderMock.AssertExpectations(t)
}
