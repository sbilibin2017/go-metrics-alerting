package repositories

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мокаем Storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockStorage) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

// Вместо <-chan [2]string используем chan [2]string
func (m *MockStorage) Generate() <-chan [2]string {
	args := m.Called()
	return args.Get(0).(chan [2]string) // Приводим к нужному типу
}

// Мокаем KeyProcessor
type MockKeyProcessor struct {
	mock.Mock
}

func (m *MockKeyProcessor) Encode(metricType string, metricName string) string {
	args := m.Called(metricType, metricName)
	return args.String(0)
}

func (m *MockKeyProcessor) Decode(key string) (string, string, error) {
	args := m.Called(key)
	return args.String(0), args.String(1), args.Error(2)
}

// Тестирование Save
func TestSave(t *testing.T) {
	mockStorage := new(MockStorage)
	mockKeyProcessor := new(MockKeyProcessor)
	repo := NewMetricRepository(mockStorage, mockKeyProcessor)

	// Мокаем поведение keyProcessor и storage
	mockKeyProcessor.On("Encode", "gauge", "metric1").Return("gauge_metric1")
	mockStorage.On("Set", "gauge_metric1", "100").Return(nil)

	// Сохраняем метрику
	err := repo.Save("gauge", "metric1", "100")

	// Проверяем, что не возникло ошибки
	assert.NoError(t, err)

	// Проверяем, что вызовы были сделаны с правильными параметрами
	mockKeyProcessor.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

// Тестирование Get
func TestGet(t *testing.T) {
	mockStorage := new(MockStorage)
	mockKeyProcessor := new(MockKeyProcessor)
	repo := NewMetricRepository(mockStorage, mockKeyProcessor)

	// Мокаем поведение keyProcessor и storage
	mockKeyProcessor.On("Encode", "gauge", "metric1").Return("gauge_metric1")
	mockStorage.On("Get", "gauge_metric1").Return("100", nil)

	// Извлекаем метрику
	value, err := repo.Get("gauge", "metric1")

	// Проверяем, что ошибка отсутствует и значение правильно
	assert.NoError(t, err)
	assert.Equal(t, "100", value)

	// Проверяем, что вызовы были сделаны с правильными параметрами
	mockKeyProcessor.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

// Тестирование Get с ошибкой (метрика не найдена)
func TestGetNotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	mockKeyProcessor := new(MockKeyProcessor)
	repo := NewMetricRepository(mockStorage, mockKeyProcessor)

	// Мокаем поведение keyProcessor и storage
	mockKeyProcessor.On("Encode", "gauge", "metric1").Return("gauge_metric1")
	mockStorage.On("Get", "gauge_metric1").Return("", errors.New("value not found"))

	// Извлекаем метрику
	_, err := repo.Get("gauge", "metric1")

	// Проверяем, что возникла ошибка
	assert.Error(t, err)

	// Проверяем, что вызовы были сделаны с правильными параметрами
	mockKeyProcessor.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

// Тестирование GetAll
func TestGetAll(t *testing.T) {
	mockStorage := new(MockStorage)
	mockKeyProcessor := new(MockKeyProcessor)
	repo := NewMetricRepository(mockStorage, mockKeyProcessor)

	// Мокаем генератор метрик и поведение keyProcessor
	metricChannel := make(chan [2]string, 2)
	metricChannel <- [2]string{"gauge_metric1", "100"}
	metricChannel <- [2]string{"counter_metric2", "200"}
	close(metricChannel)

	mockStorage.On("Generate").Return(metricChannel)
	mockKeyProcessor.On("Decode", "gauge_metric1").Return("gauge", "metric1", nil)
	mockKeyProcessor.On("Decode", "counter_metric2").Return("counter", "metric2", nil)

	// Получаем все метрики
	allMetrics := repo.GetAll()

	// Проверяем, что метрики правильно извлечены
	assert.Len(t, allMetrics, 2)
	assert.Equal(t, [3]string{"gauge", "metric1", "100"}, allMetrics[0])
	assert.Equal(t, [3]string{"counter", "metric2", "200"}, allMetrics[1])

	// Проверяем, что вызовы были сделаны с правильными параметрами
	mockStorage.AssertExpectations(t)
	mockKeyProcessor.AssertExpectations(t)
}
