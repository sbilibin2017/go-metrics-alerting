package services

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockValidateEmptyString is a mock for ValidateEmptyString
type MockIDValidator struct {
	mock.Mock
}

func (m *MockIDValidator) Validate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockValidateMType is a mock for ValidateMType
type MockValidateMType struct {
	mock.Mock
}

func (m *MockValidateMType) Validate(mType types.MType) error {
	args := m.Called(mType)
	return args.Error(0)
}

// MockValidateDelta is a mock for ValidateDelta
type MockValidateDelta struct {
	mock.Mock
}

func (m *MockValidateDelta) Validate(mType types.MType, delta *int64) error {
	args := m.Called(mType, delta)
	return args.Error(0)
}

// MockValidateValue is a mock for ValidateValue
type MockValidateValue struct {
	mock.Mock
}

func (m *MockValidateValue) Validate(mType types.MType, value *float64) error {
	args := m.Called(mType, value)
	return args.Error(0)
}

// MockValidateCounterValue is a mock for ValidateCounterValue
type MockValidateCounterValue struct {
	mock.Mock
}

func (m *MockValidateCounterValue) Validate(value string) error {
	args := m.Called(value)
	return args.Error(0)
}

// MockValidateGaugeValue is a mock for ValidateGaugeValue
type MockValidateGaugeValue struct {
	mock.Mock
}

func (m *MockValidateGaugeValue) Validate(value string) error {
	args := m.Called(value)
	return args.Error(0)
}

// MockInt64ParserFormatter is a mock for Int64ParserFormatter
type MockInt64Formatter struct {
	mock.Mock
}

func (m *MockInt64Formatter) Parse(value string) (int64, error) {
	args := m.Called(value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInt64Formatter) Format(value int64) string {
	args := m.Called(value)
	return args.String(0)
}

// MockFloat64ParserFormatter is a mock for Float64ParserFormatter
type MockFloat64Formatter struct {
	mock.Mock
}

func (m *MockFloat64Formatter) Parse(value string) (float64, error) {
	args := m.Called(value)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockFloat64Formatter) Format(value float64) string {
	args := m.Called(value)
	return args.String(0)
}

// MockSaver is a mock of the SaverInterface
type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(key, value string) bool {
	args := m.Called(key, value)
	return args.Bool(0)
}

// MockGetter is a mock of the GetterInterface
type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

// MockRanger is a mock of the RangerInterface
type MockRanger struct {
	mock.Mock
}

func (m *MockRanger) Range(callback func(key, value string) bool) {
	m.Called(callback)
}

func TestUpdateMetric_IDValidation_Failure(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "invalidID").Return(assert.AnError) // Return an error for the invalid ID

	// Mock other dependencies (they are not used in this test)
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockMTypeValidator := new(MockValidateMType)
	mockDeltaValidator := new(MockValidateDelta)
	mockValueValidator := new(MockValidateValue)
	mockInt64Formatter := new(MockInt64Formatter)
	mockFloat64Formatter := new(MockFloat64Formatter)

	// Create the UpdateMetricsService with mocked dependencies
	updateService := NewUpdateMetricsService(
		mockSaver,
		mockGetter,
		mockIDValidator,
		mockMTypeValidator,
		mockDeltaValidator,
		mockValueValidator,
		mockInt64Formatter,
		mockFloat64Formatter,
	)

	// Create a request with an invalid ID
	req := &types.UpdateMetricsRequest{
		ID: "invalidID",
		// Other fields like MType, Delta, Value can be left empty as this test is only for ID validation
	}

	// Act
	resp, errResp := updateService.UpdateMetricValue(req)

	// Assert
	// We expect an error response with status 404 and the appropriate message
	assert.Nil(t, resp)       // No successful response
	assert.NotNil(t, errResp) // Error response should be present
	assert.Equal(t, http.StatusNotFound, errResp.Status)
	assert.Equal(t, "Metric with the given ID not found", errResp.Message)

	// Assert that the Validate method was called with the correct ID
	mockIDValidator.AssertCalled(t, "Validate", "invalidID")
}

func TestUpdateMetric_MTypeValidation_Failure(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // ID валиден

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.MType("invalidType")).Return(assert.AnError) // Возвращаем ошибку при валидации типа метрики

	// Мокаем остальные зависимости
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockDeltaValidator := new(MockValidateDelta)
	mockValueValidator := new(MockValidateValue)
	mockInt64Formatter := new(MockInt64Formatter)
	mockFloat64Formatter := new(MockFloat64Formatter)

	// Создаем UpdateMetricsService с моками
	updateService := NewUpdateMetricsService(
		mockSaver,
		mockGetter,
		mockIDValidator,
		mockMTypeValidator,
		mockDeltaValidator,
		mockValueValidator,
		mockInt64Formatter,
		mockFloat64Formatter,
	)

	// Создаем запрос с некорректным MType
	req := &types.UpdateMetricsRequest{
		ID:    "validID",                  // Валидный ID
		MType: types.MType("invalidType"), // Невалидный тип метрики
		// Остальные поля оставляем пустыми, так как в этом тесте проверяем только MType
	}

	// Act
	resp, errResp := updateService.UpdateMetricValue(req)

	// Assert
	// Ожидаем ошибку с кодом 400 и соответствующим сообщением
	assert.Nil(t, resp)       // Ответ не должен быть успешным
	assert.NotNil(t, errResp) // Ошибка должна быть
	assert.Equal(t, http.StatusBadRequest, errResp.Status)
	assert.Equal(t, "Invalid metric type", errResp.Message)

	// Убеждаемся, что метод Validate был вызван с правильным MType
	mockMTypeValidator.AssertCalled(t, "Validate", types.MType("invalidType"))
}

func TestUpdateMetric_DeltaValidation_Failure(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // ID валиден

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Counter).Return(nil) // Валидный тип метрики (Counter)

	mockDeltaValidator := new(MockValidateDelta)
	mockDeltaValidator.On("Validate", types.Counter, mock.Anything).Return(assert.AnError) // Возвращаем ошибку для Delta в случае типа Counter

	// Мокаем остальные зависимости
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockValueValidator := new(MockValidateValue)
	mockInt64Formatter := new(MockInt64Formatter)
	mockFloat64Formatter := new(MockFloat64Formatter)

	// Создаем UpdateMetricsService с моками
	updateService := NewUpdateMetricsService(
		mockSaver,
		mockGetter,
		mockIDValidator,
		mockMTypeValidator,
		mockDeltaValidator,
		mockValueValidator,
		mockInt64Formatter,
		mockFloat64Formatter,
	)

	// Создаем запрос с некорректным Delta значением
	req := &types.UpdateMetricsRequest{
		ID:    "validID",     // Валидный ID
		MType: types.Counter, // Тип метрики Counter
		Delta: new(int64),    // Некорректное значение для Delta (например, нулевое)
		// Остальные поля оставляем пустыми, так как в этом тесте проверяем только Delta
	}

	// Act
	resp, errResp := updateService.UpdateMetricValue(req)

	// Assert
	// Ожидаем ошибку с кодом 400 и соответствующим сообщением
	assert.Nil(t, resp)       // Ответ не должен быть успешным
	assert.NotNil(t, errResp) // Ошибка должна быть
	assert.Equal(t, http.StatusBadRequest, errResp.Status)
	assert.Equal(t, "Invalid delta value", errResp.Message)

	// Убеждаемся, что метод Validate был вызван с правильным MType и Delta
	mockDeltaValidator.AssertCalled(t, "Validate", types.Counter, mock.Anything)
}

func TestUpdateMetric_ValueValidation_Failure(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // Валидный ID

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Gauge).Return(nil) // Валидный тип метрики (Gauge)

	mockDeltaValidator := new(MockValidateDelta)
	mockDeltaValidator.On("Validate", types.Gauge, mock.Anything).Return(nil) // Валидный Delta

	mockValueValidator := new(MockValidateValue)
	mockValueValidator.On("Validate", types.Gauge, mock.Anything).Return(assert.AnError) // Возвращаем ошибку для значения Value

	// Мокаем остальные зависимости
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockInt64Formatter := new(MockInt64Formatter)
	mockFloat64Formatter := new(MockFloat64Formatter)

	// Создаем UpdateMetricsService с моками
	updateService := NewUpdateMetricsService(
		mockSaver,
		mockGetter,
		mockIDValidator,
		mockMTypeValidator,
		mockDeltaValidator,
		mockValueValidator,
		mockInt64Formatter,
		mockFloat64Formatter,
	)

	// Создаем запрос с некорректным значением Value
	req := &types.UpdateMetricsRequest{
		ID:    "validID",    // Валидный ID
		MType: types.Gauge,  // Тип метрики Gauge
		Value: new(float64), // Некорректное значение для Value (например, nil или ошибочное значение)
		// Остальные поля оставляем пустыми
	}

	// Act
	resp, errResp := updateService.UpdateMetricValue(req)

	// Assert
	// Ожидаем ошибку с кодом 400 и соответствующим сообщением
	assert.Nil(t, resp)       // Ответ не должен быть успешным
	assert.NotNil(t, errResp) // Ошибка должна быть
	assert.Equal(t, http.StatusBadRequest, errResp.Status)
	assert.Equal(t, "Invalid value", errResp.Message)

	// Убеждаемся, что метод Validate был вызван с правильным MType и Value
	mockValueValidator.AssertCalled(t, "Validate", types.Gauge, mock.Anything)
}

func TestUpdateMetric_GaugeSave_Success(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // Валидный ID

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Gauge).Return(nil) // Валидный тип метрики (Gauge)

	mockDeltaValidator := new(MockValidateDelta)
	mockDeltaValidator.On("Validate", types.Gauge, mock.Anything).Return(nil) // Валидный Delta

	mockValueValidator := new(MockValidateValue)
	mockValueValidator.On("Validate", types.Gauge, mock.Anything).Return(nil) // Валидная валидация значения

	// Мокаем зависимости
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockInt64Formatter := new(MockInt64Formatter)
	mockFloat64Formatter := new(MockFloat64Formatter)

	// Настроим моки для Float64Formatter
	mockFloat64Formatter.On("Format", float64(3.14)).Return("3.14") // Форматируем значение 3.14 как "3.14"

	// Настроим мок для Save
	mockSaver.On("Save", "validID", "3.14").Return(true) // Ожидаем вызов Save с параметрами "validID" и "3.14"

	// Создаем UpdateMetricsService с моками
	updateService := NewUpdateMetricsService(
		mockSaver,
		mockGetter,
		mockIDValidator,
		mockMTypeValidator,
		mockDeltaValidator,
		mockValueValidator,
		mockInt64Formatter,
		mockFloat64Formatter,
	)

	// Создаем запрос с валидными данными
	req := &types.UpdateMetricsRequest{
		ID:    "validID",    // Валидный ID
		MType: types.Gauge,  // Тип метрики Gauge
		Value: new(float64), // Значение, которое мы хотим сохранить
	}
	*req.Value = 3.14 // Устанавливаем значение метрики

	// Act
	resp, errResp := updateService.UpdateMetricValue(req)

	// Assert
	// Проверяем, что ошибки нет
	assert.Nil(t, errResp)
	assert.NotNil(t, resp)

	// Проверяем, что метод Save был вызван с правильными параметрами
	mockSaver.AssertCalled(t, "Save", "validID", "3.14")

	// Проверяем, что форматирование значения было вызвано
	mockFloat64Formatter.AssertCalled(t, "Format", float64(3.14))
}

func TestUpdateMetric_CounterSave_Success(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // Валидный ID

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Counter).Return(nil) // Валидный тип метрики (Counter)

	mockDeltaValidator := new(MockValidateDelta)
	mockDeltaValidator.On("Validate", types.Counter, mock.Anything).Return(nil) // Валидный Delta

	mockValueValidator := new(MockValidateValue)
	mockValueValidator.On("Validate", types.Counter, mock.Anything).Return(nil) // Валидная валидация значения

	// Мокаем зависимости
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockInt64Formatter := new(MockInt64Formatter)
	mockFloat64Formatter := new(MockFloat64Formatter)

	// Мокаем поведение для Getter
	mockGetter.On("Get", "validID").Return("5", true) // Возвращаем значение "5" для существующего ключа

	// Мокаем поведение для Int64Formatter
	mockInt64Formatter.On("Parse", "5").Return(int64(5), nil) // Парсим строку "5" как int64(5)
	mockInt64Formatter.On("Format", int64(10)).Return("10")   // Форматируем число 10 как строку "10"

	// Настроим мок для Save
	mockSaver.On("Save", "validID", "10").Return(true) // Ожидаем вызов Save с параметрами "validID" и "10"

	// Создаем UpdateMetricsService с моками
	updateService := NewUpdateMetricsService(
		mockSaver,
		mockGetter,
		mockIDValidator,
		mockMTypeValidator,
		mockDeltaValidator,
		mockValueValidator,
		mockInt64Formatter,
		mockFloat64Formatter,
	)

	// Создаем запрос с валидными данными
	req := &types.UpdateMetricsRequest{
		ID:    "validID",     // Валидный ID
		MType: types.Counter, // Тип метрики Counter
		Delta: new(int64),    // Значение Delta
	}
	*req.Delta = 5 // Устанавливаем Delta равным 5

	// Act
	resp, errResp := updateService.UpdateMetricValue(req)

	// Assert
	// Проверяем, что ошибки нет
	assert.Nil(t, errResp)
	assert.NotNil(t, resp)

	// Проверяем, что метод Save был вызван с правильными параметрами
	mockSaver.AssertCalled(t, "Save", "validID", "10") // Ожидаем вызов Save с ID и новым значением

	// Проверяем, что метод Get был вызван с правильным ID
	mockGetter.AssertCalled(t, "Get", "validID")

	// Проверяем, что парсинг значения был вызван с правильной строкой
	mockInt64Formatter.AssertCalled(t, "Parse", "5")

	// Проверяем, что форматирование нового значения Delta было вызвано
	mockInt64Formatter.AssertCalled(t, "Format", int64(10))
}

func TestParseMetricValues_CounterValid(t *testing.T) {
	// Arrange
	mockInt64Formatter := new(MockInt64Formatter)
	mockInt64Formatter.On("Parse", "123").Return(int64(123), nil) // Возвращаем корректное значение

	// Создаем сервис с моками
	service := &UpdateMetricsService{
		Int64Formatter: mockInt64Formatter,
	}

	// Act
	value, delta, errResp := service.ParseMetricValues(string(types.Counter), "123")

	// Assert
	assert.Nil(t, errResp)              // Ошибки не ожидается
	assert.Equal(t, int64(123), *delta) // Проверяем, что delta имеет правильное значение
	assert.Nil(t, value)                // Значение для Counter должно быть nil
}

func TestParseMetricValues_CounterInvalidValue(t *testing.T) {
	// Arrange
	mockInt64Formatter := new(MockInt64Formatter)
	mockInt64Formatter.On("Parse", "abc").Return(int64(0), assert.AnError) // Некорректное значение

	// Создаем сервис с моками
	service := &UpdateMetricsService{
		Int64Formatter: mockInt64Formatter,
	}

	// Act
	value, delta, errResp := service.ParseMetricValues(string(types.Counter), "abc")

	// Assert
	assert.NotNil(t, errResp)                                            // Ошибка ожидается
	assert.Equal(t, http.StatusBadRequest, errResp.Status)               // Ошибка 400
	assert.Equal(t, "Invalid value format for Counter", errResp.Message) // Сообщение об ошибке
	assert.Nil(t, value)                                                 // Значение должно быть nil
	assert.Nil(t, delta)                                                 // Delta должно быть nil
}

func TestParseMetricValues_GaugeValid(t *testing.T) {
	// Arrange
	mockFloat64Formatter := new(MockFloat64Formatter)
	mockFloat64Formatter.On("Parse", "3.14").Return(float64(3.14), nil) // Возвращаем корректное значение

	// Создаем сервис с моками
	service := &UpdateMetricsService{
		Float64Formatter: mockFloat64Formatter,
	}

	// Act
	value, delta, errResp := service.ParseMetricValues(string(types.Gauge), "3.14")

	// Assert
	assert.Nil(t, errResp)                 // Ошибки не ожидается
	assert.Equal(t, float64(3.14), *value) // Проверяем, что value имеет правильное значение
	assert.Nil(t, delta)                   // Delta для Gauge должно быть nil
}

func TestParseMetricValues_GaugeInvalidValue(t *testing.T) {
	// Arrange
	mockFloat64Formatter := new(MockFloat64Formatter)
	mockFloat64Formatter.On("Parse", "abc").Return(float64(0), assert.AnError) // Некорректное значение

	// Создаем сервис с моками
	service := &UpdateMetricsService{
		Float64Formatter: mockFloat64Formatter,
	}

	// Act
	value, delta, errResp := service.ParseMetricValues(string(types.Gauge), "abc")

	// Assert
	assert.NotNil(t, errResp)                                          // Ошибка ожидается
	assert.Equal(t, http.StatusBadRequest, errResp.Status)             // Ошибка 400
	assert.Equal(t, "Invalid value format for Gauge", errResp.Message) // Сообщение об ошибке
	assert.Nil(t, value)                                               // Значение должно быть nil
	assert.Nil(t, delta)                                               // Delta должно быть nil
}

func TestParseMetricValues_UnknownMType(t *testing.T) {
	// Arrange
	service := &UpdateMetricsService{}

	// Act
	value, delta, errResp := service.ParseMetricValues("unknownMetricType", "3.14")

	// Assert
	assert.Nil(t, errResp) // Ошибки не ожидается
	assert.Nil(t, value)   // Не должно быть значения для неизвестного типа
	assert.Nil(t, delta)   // Не должно быть delta для неизвестного типа
}

func TestGetMetricValue_IDValidation_Failure(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "invalidID").Return(assert.AnError) // Ошибка при валидации ID

	mockGetter := new(MockGetter)
	mockMTypeValidator := new(MockValidateMType)

	// Создаем сервис с моками
	service := NewGetMetricValueService(mockGetter, mockIDValidator, mockMTypeValidator)

	// Act
	req := &types.GetMetricValueRequest{
		ID:    "invalidID",
		MType: types.Gauge, // Тип метрики не имеет значения для этого теста
	}
	resp, errResp := service.GetMetricValue(req)

	// Assert
	assert.Nil(t, resp)                                  // Ответ должен быть nil
	assert.NotNil(t, errResp)                            // Ошибка должна быть
	assert.Equal(t, http.StatusNotFound, errResp.Status) // Ожидаем ошибку 404
	assert.Equal(t, "Metric with the given ID not found", errResp.Message)
}

func TestGetMetricValue_MTypeValidation_Failure(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // Валидация ID проходит

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Gauge).Return(assert.AnError) // Ошибка при валидации типа метрики

	mockGetter := new(MockGetter)

	// Создаем сервис с моками
	service := NewGetMetricValueService(mockGetter, mockIDValidator, mockMTypeValidator)

	// Act
	req := &types.GetMetricValueRequest{
		ID:    "validID",
		MType: types.Gauge, // Используем тип Gauge
	}
	resp, errResp := service.GetMetricValue(req)

	// Assert
	assert.Nil(t, resp)                                    // Ответ должен быть nil
	assert.NotNil(t, errResp)                              // Ошибка должна быть
	assert.Equal(t, http.StatusBadRequest, errResp.Status) // Ожидаем ошибку 400
	assert.Equal(t, "Invalid metric type", errResp.Message)
}

func TestGetMetricValue_MetricNotFound(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // Валидация ID проходит

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Gauge).Return(nil) // Валидация типа метрики проходит

	mockGetter := new(MockGetter)
	mockGetter.On("Get", "validID").Return("", false) // Не находим метрику по ID

	// Создаем сервис с моками
	service := NewGetMetricValueService(mockGetter, mockIDValidator, mockMTypeValidator)

	// Act
	req := &types.GetMetricValueRequest{
		ID:    "validID",
		MType: types.Gauge,
	}
	resp, errResp := service.GetMetricValue(req)

	// Assert
	assert.Nil(t, resp)                                  // Ответ должен быть nil
	assert.NotNil(t, errResp)                            // Ошибка должна быть
	assert.Equal(t, http.StatusNotFound, errResp.Status) // Ожидаем ошибку 404
	assert.Equal(t, "Metric not found", errResp.Message)
}

func TestGetMetricValue_Success(t *testing.T) {
	// Arrange
	mockIDValidator := new(MockIDValidator)
	mockIDValidator.On("Validate", "validID").Return(nil) // Валидация ID проходит

	mockMTypeValidator := new(MockValidateMType)
	mockMTypeValidator.On("Validate", types.Gauge).Return(nil) // Валидация типа метрики проходит

	mockGetter := new(MockGetter)
	mockGetter.On("Get", "validID").Return("3.14", true) // Находим метрику с значением "3.14"

	// Создаем сервис с моками
	service := NewGetMetricValueService(mockGetter, mockIDValidator, mockMTypeValidator)

	// Act
	req := &types.GetMetricValueRequest{
		ID:    "validID",
		MType: types.Gauge,
	}
	resp, errResp := service.GetMetricValue(req)

	// Assert
	assert.Nil(t, errResp)              // Ошибки не ожидается
	assert.NotNil(t, resp)              // Ответ должен быть
	assert.Equal(t, "validID", resp.ID) // Проверяем ID
	assert.Equal(t, "3.14", resp.Value) // Проверяем значение метрики
}

func TestGetAllMetricValues_Success(t *testing.T) {
	// Arrange
	mockRanger := new(MockRanger)

	// Подготовим данные, которые Ranger должен возвращать
	mockRanger.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(key, value string) bool)
		fn("metric1", "3.14")
		fn("metric2", "2.71")
		fn("metric3", "1.62")
	}).Once()

	// Создаем сервис с моками
	service := NewGetAllMetricValuesService(mockRanger)

	// Act
	metrics, errResp := service.GetAllMetricValues()

	// Assert
	assert.Nil(t, errResp)                    // Ошибок не ожидаем
	assert.Len(t, metrics, 3)                 // Должно быть 3 метрики
	assert.Equal(t, "metric1", metrics[0].ID) // Проверка на правильность ID
	assert.Equal(t, "3.14", metrics[0].Value) // Проверка на правильность значения
	assert.Equal(t, "metric2", metrics[1].ID)
	assert.Equal(t, "2.71", metrics[1].Value)
	assert.Equal(t, "metric3", metrics[2].ID)
	assert.Equal(t, "1.62", metrics[2].Value)
}

func TestGetAllMetricValues_SingleMetric(t *testing.T) {
	// Arrange
	mockRanger := new(MockRanger)

	// Подготовим одну метрику
	mockRanger.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(key, value string) bool)
		fn("metric1", "3.14")
	}).Once()

	// Создаем сервис с моками
	service := NewGetAllMetricValuesService(mockRanger)

	// Act
	metrics, errResp := service.GetAllMetricValues()

	// Assert
	assert.Nil(t, errResp)                    // Ошибок не ожидаем
	assert.Len(t, metrics, 1)                 // Должна быть одна метрика
	assert.Equal(t, "metric1", metrics[0].ID) // Проверка на правильность ID
	assert.Equal(t, "3.14", metrics[0].Value) // Проверка на правильность значения
}
