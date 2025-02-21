package types

import (
	"errors"
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мокируем валидаторы

type MockStringValidator struct {
	mock.Mock
}

func (m *MockStringValidator) Validate(s string) error {
	args := m.Called(s)
	return args.Error(0)
}

type MockMetricTypeValidator struct {
	mock.Mock
}

func (m *MockMetricTypeValidator) Validate(mType domain.MType) error {
	args := m.Called(mType)
	return args.Error(0)
}

type MockDeltaValidator struct {
	mock.Mock
}

func (m *MockDeltaValidator) Validate(mType domain.MType, delta *int64) error {
	args := m.Called(mType, delta)
	return args.Error(0)
}

type MockValueValidator struct {
	mock.Mock
}

func (m *MockValueValidator) Validate(mType domain.MType, value *float64) error {
	args := m.Called(mType, value)
	return args.Error(0)
}

type MockValueStringValidator struct {
	mock.Mock
}

func (m *MockValueStringValidator) Validate(value string) error {
	args := m.Called(value)
	return args.Error(0)
}

// Тесты

func TestUpdateMetricBodyRequest_Validate_Success(t *testing.T) {
	// Мокируем валидаторы
	stringValidator := new(MockStringValidator)
	metricTypeValidator := new(MockMetricTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	// Создаем тестовые данные
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "counter",
		Delta: new(int64),
	}

	// Ожидаем, что все валидаторы будут успешно пройдены
	stringValidator.On("Validate", request.ID).Return(nil)
	metricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	deltaValidator.On("Validate", domain.MType(request.MType), request.Delta).Return(nil)

	// Запускаем тест
	err := request.Validate(stringValidator, metricTypeValidator, deltaValidator, valueValidator)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что все ожидания выполнены
	stringValidator.AssertExpectations(t)
	metricTypeValidator.AssertExpectations(t)
	deltaValidator.AssertExpectations(t)
}

func TestUpdateMetricBodyRequest_Validate_Failure_InvalidID(t *testing.T) {
	// Создаем структуру запроса с некорректным ID
	request := &UpdateMetricBodyRequest{
		ID:    "", // Пустой ID
		MType: "counter",
		Delta: new(int64),
		Value: nil,
	}

	// Устанавливаем значение для Delta
	*request.Delta = 10

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockDeltaValidator := new(MockDeltaValidator)
	mockValueValidator := new(MockValueValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(errors.New("invalid ID"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockDeltaValidator, mockValueValidator)

	// Проверяем, что ошибка возникла из-за некорректного ID
	assert.Error(t, err)
	assert.Equal(t, "invalid ID", err.Error())

	// Проверяем, что остальные методы не были вызваны, так как валидация по ID не прошла
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertNotCalled(t, "Validate")
	mockDeltaValidator.AssertNotCalled(t, "Validate")
	mockValueValidator.AssertNotCalled(t, "Validate")
}

func TestUpdateMetricBodyRequest_Validate_Failure_InvalidMetricType(t *testing.T) {
	// Создаем структуру запроса с некорректным типом метрики
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "unknown", // Некорректный тип
		Delta: new(int64),
		Value: nil,
	}

	// Устанавливаем значение для Delta
	*request.Delta = 10

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockDeltaValidator := new(MockDeltaValidator)
	mockValueValidator := new(MockValueValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(errors.New("invalid metric type"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockDeltaValidator, mockValueValidator)

	// Проверяем, что ошибка возникла из-за некорректного типа метрики
	assert.Error(t, err)
	assert.Equal(t, "invalid metric type", err.Error())

	// Проверяем, что другие валидаторы не были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockDeltaValidator.AssertNotCalled(t, "Validate")
	mockValueValidator.AssertNotCalled(t, "Validate")
}

func TestUpdateMetricBodyRequest_Validate_Success_Counter(t *testing.T) {
	// Создаем структуру запроса с корректными значениями
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "counter",
		Delta: new(int64),
		Value: nil,
	}

	// Устанавливаем значение для Delta
	*request.Delta = 10

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockDeltaValidator := new(MockDeltaValidator)
	mockValueValidator := new(MockValueValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	mockDeltaValidator.On("Validate", domain.MType(request.MType), request.Delta).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockDeltaValidator, mockValueValidator)

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что все моковые методы были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockDeltaValidator.AssertExpectations(t)
	mockValueValidator.AssertNotCalled(t, "Validate")
}

func TestUpdateMetricBodyRequest_Validate_Success_Gauge(t *testing.T) {
	// Создаем структуру запроса с корректными значениями
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "gauge",
		Delta: nil,
		Value: new(float64),
	}

	// Устанавливаем значение для Value
	*request.Value = 10.5

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockDeltaValidator := new(MockDeltaValidator)
	mockValueValidator := new(MockValueValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	mockValueValidator.On("Validate", domain.MType(request.MType), request.Value).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockDeltaValidator, mockValueValidator)

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что все моковые методы были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockDeltaValidator.AssertNotCalled(t, "Validate")
	mockValueValidator.AssertExpectations(t)
}

func TestUpdateMetricBodyRequest_Validate_Failure_Counter_InvalidDelta(t *testing.T) {
	// Создаем структуру запроса с некорректным значением Delta для типа Counter
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "counter",
		Delta: new(int64),
		Value: nil,
	}

	// Устанавливаем некорректное значение для Delta
	*request.Delta = 0 // Здесь значение для Delta не подходит

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockDeltaValidator := new(MockDeltaValidator)
	mockValueValidator := new(MockValueValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	mockDeltaValidator.On("Validate", domain.MType(request.MType), request.Delta).Return(errors.New("invalid delta"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockDeltaValidator, mockValueValidator)

	// Проверяем, что ошибка вернулась
	assert.Error(t, err)
	assert.Equal(t, "invalid delta", err.Error())

	// Проверяем, что другие валидаторы не были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockDeltaValidator.AssertExpectations(t)
	mockValueValidator.AssertNotCalled(t, "Validate")
}

func TestUpdateMetricBodyRequest_Validate_Failure_Gauge_InvalidValue(t *testing.T) {
	// Создаем структуру запроса с некорректным значением Value для типа Gauge
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "gauge",
		Delta: nil,
		Value: new(float64),
	}

	// Устанавливаем некорректное значение для Value
	*request.Value = -10.5 // Здесь значение для Value не подходит

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockDeltaValidator := new(MockDeltaValidator)
	mockValueValidator := new(MockValueValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	mockValueValidator.On("Validate", domain.MType(request.MType), request.Value).Return(errors.New("invalid value"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockDeltaValidator, mockValueValidator)

	// Проверяем, что ошибка вернулась
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())

	// Проверяем, что другие валидаторы не были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockDeltaValidator.AssertNotCalled(t, "Validate")
	mockValueValidator.AssertExpectations(t)
}

func TestUpdateMetricBodyRequest_Validate_Failure(t *testing.T) {
	// Мокируем валидаторы
	stringValidator := new(MockStringValidator)
	metricTypeValidator := new(MockMetricTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	// Создаем тестовые данные
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "counter",
		Delta: new(int64),
	}

	// Ожидаем, что валидатор строки пройдет успешно
	stringValidator.On("Validate", request.ID).Return(nil)
	// Ожидаем, что валидатор типа метрики вернет ошибку
	metricTypeValidator.On("Validate", domain.MType(request.MType)).Return(errors.New("invalid metric type"))
	// Не ожидаем вызова deltaValidator, так как ошибка уже будет на этапе проверки типа метрики

	// Запускаем тест
	err := request.Validate(stringValidator, metricTypeValidator, deltaValidator, valueValidator)

	// Проверяем, что ошибка возникла
	assert.Error(t, err)
	assert.Equal(t, "invalid metric type", err.Error())

	// Проверяем, что все ожидания выполнены
	stringValidator.AssertExpectations(t)
	metricTypeValidator.AssertExpectations(t)
	deltaValidator.AssertExpectations(t) // Проверяем, что deltaValidator не был вызван
}

func TestUpdateMetricPathRequest_Validate_Success(t *testing.T) {
	// Мокируем валидаторы
	stringValidator := new(MockStringValidator)
	metricTypeValidator := new(MockMetricTypeValidator)
	valueStringValidator := new(MockValueStringValidator)

	// Создаем тестовые данные
	request := &UpdateMetricPathRequest{
		ID:    "metric2",
		MType: "gauge",
		Value: "10.5",
	}

	// Ожидаем, что все валидаторы будут успешно пройдены
	stringValidator.On("Validate", request.ID).Return(nil)
	metricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	valueStringValidator.On("Validate", request.Value).Return(nil)

	// Запускаем тест
	err := request.Validate(stringValidator, metricTypeValidator, valueStringValidator)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что все ожидания выполнены
	stringValidator.AssertExpectations(t)
	metricTypeValidator.AssertExpectations(t)
	valueStringValidator.AssertExpectations(t)
}

func TestUpdateMetricPathRequest_Validate_Failure(t *testing.T) {
	// Мокируем валидаторы
	stringValidator := new(MockStringValidator)
	metricTypeValidator := new(MockMetricTypeValidator)
	valueStringValidator := new(MockValueStringValidator)

	// Создаем тестовые данные
	request := &UpdateMetricPathRequest{
		ID:    "metric2",
		MType: "gauge",
		Value: "invalid",
	}

	// Ожидаем, что валидатор значения метрики вернет ошибку
	stringValidator.On("Validate", request.ID).Return(nil)
	metricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	valueStringValidator.On("Validate", request.Value).Return(errors.New("invalid value"))

	// Запускаем тест
	err := request.Validate(stringValidator, metricTypeValidator, valueStringValidator)

	// Проверяем, что ошибка возникла
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())

	// Проверяем, что все ожидания выполнены
	stringValidator.AssertExpectations(t)
	metricTypeValidator.AssertExpectations(t)
	valueStringValidator.AssertExpectations(t)
}

func TestUpdateMetricPathRequest_Validate_Failure_InvalidID(t *testing.T) {
	// Создаем структуру запроса с некорректным значением ID
	request := &UpdateMetricPathRequest{
		ID:    "",
		MType: "counter",
		Value: "100",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockValueStringValidator := new(MockValueStringValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(errors.New("invalid ID"))
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	mockValueStringValidator.On("Validate", request.Value).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockValueStringValidator)

	// Проверяем, что ошибка для ID возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid ID", err.Error())

	// Проверяем, что остальные валидаторы не были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertNotCalled(t, "Validate")
	mockValueStringValidator.AssertNotCalled(t, "Validate")
}

func TestUpdateMetricPathRequest_Validate_Failure_InvalidMType(t *testing.T) {
	// Создаем структуру запроса с некорректным значением MType
	request := &UpdateMetricPathRequest{
		ID:    "metric1",
		MType: "unknown",
		Value: "100",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockValueStringValidator := new(MockValueStringValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(errors.New("invalid MType"))
	mockValueStringValidator.On("Validate", request.Value).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockValueStringValidator)

	// Проверяем, что ошибка для MType возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid MType", err.Error())

	// Проверяем, что остальные валидаторы не были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockValueStringValidator.AssertNotCalled(t, "Validate")
}

func TestUpdateMetricPathRequest_Validate_Failure_InvalidValue(t *testing.T) {
	// Создаем структуру запроса с некорректным значением Value
	request := &UpdateMetricPathRequest{
		ID:    "metric1",
		MType: "counter",
		Value: "invalid_value", // Некорректное значение для counter
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)
	mockValueStringValidator := new(MockValueStringValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)
	mockValueStringValidator.On("Validate", request.Value).Return(errors.New("invalid value"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator, mockValueStringValidator)

	// Проверяем, что ошибка для Value возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())

	// Проверяем, что другие валидаторы не были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
	mockValueStringValidator.AssertExpectations(t)
}

func TestToDomain(t *testing.T) {
	// Создаем тестовые данные
	request := &UpdateMetricPathRequest{
		ID:    "metric3",
		MType: "gauge",
		Value: "20.5",
	}

	// Преобразуем в доменную сущность
	domainMetric := request.ToDomain()

	// Проверяем, что доменная сущность содержит правильные данные
	assert.Equal(t, "metric3", domainMetric.ID)
	assert.Equal(t, domain.Gauge, domainMetric.MType)
	assert.Nil(t, domainMetric.Delta)
	assert.NotNil(t, domainMetric.Value)
	assert.Equal(t, 20.5, *domainMetric.Value)
}

func TestToDomainCounter(t *testing.T) {
	// Создаем тестовые данные для типа Counter
	request := &UpdateMetricPathRequest{
		ID:    "metric4",
		MType: "counter",
		Value: "10",
	}

	// Преобразуем в доменную сущность
	domainMetric := request.ToDomain()

	// Проверяем, что доменная сущность содержит правильные данные
	assert.Equal(t, "metric4", domainMetric.ID)
	assert.Equal(t, domain.Counter, domainMetric.MType)
	assert.NotNil(t, domainMetric.Delta)
	assert.Nil(t, domainMetric.Value)
	assert.Equal(t, int64(10), *domainMetric.Delta)
}

func TestUpdateMetricBodyRequest_ToDomain(t *testing.T) {
	// Создаем структуру запроса с тестовыми данными
	request := &UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "counter",
		Delta: new(int64),
		Value: nil,
	}

	// Устанавливаем значение для Delta
	*request.Delta = 10

	// Преобразуем запрос в доменную сущность
	domainMetric := request.ToDomain()

	// Проверяем, что ID и MType преобразованы правильно
	assert.Equal(t, request.ID, domainMetric.ID)
	assert.Equal(t, domain.MType(request.MType), domainMetric.MType)
	assert.Equal(t, request.Delta, domainMetric.Delta)
	assert.Nil(t, domainMetric.Value) // Value должно быть nil, так как для Counter используется Delta
}

func TestUpdateMetricBodyRequest_ToDomain_Gauge(t *testing.T) {
	// Создаем структуру запроса для типа метрики Gauge
	request := &UpdateMetricBodyRequest{
		ID:    "metric2",
		MType: "gauge",
		Delta: nil,
		Value: new(float64),
	}

	// Устанавливаем значение для Value
	*request.Value = 25.5

	// Преобразуем запрос в доменную сущность
	domainMetric := request.ToDomain()

	// Проверяем, что ID и MType преобразованы правильно
	assert.Equal(t, request.ID, domainMetric.ID)
	assert.Equal(t, domain.MType(request.MType), domainMetric.MType)
	assert.Nil(t, domainMetric.Delta) // Delta должно быть nil, так как для Gauge используется Value
	assert.Equal(t, request.Value, domainMetric.Value)
}

func TestGetMetricBodyRequest_Validate_Success(t *testing.T) {
	// Создаем структуру запроса с корректными значениями
	request := &GetMetricBodyRequest{
		ID:    "metric1",
		MType: "gauge",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator)

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что все моковые методы были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
}

func TestGetMetricBodyRequest_Validate_Failure_InvalidID(t *testing.T) {
	// Создаем структуру запроса с некорректным значением ID
	request := &GetMetricBodyRequest{
		ID:    "",
		MType: "gauge",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(errors.New("invalid ID"))
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator)

	// Проверяем, что ошибка для ID возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid ID", err.Error())

	// Проверяем, что второй валидатор не был вызван
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertNotCalled(t, "Validate")
}

func TestGetMetricBodyRequest_Validate_Failure_InvalidMType(t *testing.T) {
	// Создаем структуру запроса с некорректным значением MType
	request := &GetMetricBodyRequest{
		ID:    "metric1",
		MType: "unknown",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(errors.New("invalid MType"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator)

	// Проверяем, что ошибка для MType возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid MType", err.Error())

	// Проверяем, что первый валидатор был вызван
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
}

func TestGetMetricPathRequest_Validate_Success(t *testing.T) {
	// Создаем структуру запроса с корректными значениями
	request := &GetMetricPathRequest{
		ID:    "metric1",
		MType: "gauge",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator)

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что все моковые методы были вызваны
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
}

func TestGetMetricPathRequest_Validate_Failure_InvalidID(t *testing.T) {
	// Создаем структуру запроса с некорректным значением ID
	request := &GetMetricPathRequest{
		ID:    "",
		MType: "gauge",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(errors.New("invalid ID"))
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(nil)

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator)

	// Проверяем, что ошибка для ID возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid ID", err.Error())

	// Проверяем, что второй валидатор не был вызван
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertNotCalled(t, "Validate")
}

func TestGetMetricPathRequest_Validate_Failure_InvalidMType(t *testing.T) {
	// Создаем структуру запроса с некорректным значением MType
	request := &GetMetricPathRequest{
		ID:    "metric1",
		MType: "unknown",
	}

	// Моки для валидаторов
	mockStringValidator := new(MockStringValidator)
	mockMetricTypeValidator := new(MockMetricTypeValidator)

	// Настроим поведение моков
	mockStringValidator.On("Validate", request.ID).Return(nil)
	mockMetricTypeValidator.On("Validate", domain.MType(request.MType)).Return(errors.New("invalid MType"))

	// Вызываем метод Validate
	err := request.Validate(mockStringValidator, mockMetricTypeValidator)

	// Проверяем, что ошибка для MType возвращена
	assert.Error(t, err)
	assert.Equal(t, "invalid MType", err.Error())

	// Проверяем, что первый валидатор был вызван
	mockStringValidator.AssertExpectations(t)
	mockMetricTypeValidator.AssertExpectations(t)
}
