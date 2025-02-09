package repositories

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorageEngine - мок для StorageEngine
type MockStorageEngine struct {
	mock.Mock
}

func (m *MockStorageEngine) Set(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockStorageEngine) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockStorageEngine) Generate(ctx context.Context) <-chan []string {
	args := m.Called(ctx)
	return args.Get(0).(chan []string)
}

// MockKeyEngine - мок для KeyEngine
type MockKeyEngine struct {
	mock.Mock
}

func (m *MockKeyEngine) Encode(mt, mn string) string {
	args := m.Called(mt, mn)
	return args.String(0)
}

func (m *MockKeyEngine) Decode(key string) (string, string, error) {
	args := m.Called(key)
	return args.String(0), args.String(1), args.Error(2)
}

func TestMetricRepository_Save_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	repo := &MetricRepository{StorageEngine: mockStorage, KeyEngine: mockKeyEngine}

	metricType := "counter"
	metricName := "requests"
	metricValue := "100"

	// Mocking Encode to return the correct key
	mockKeyEngine.On("Encode", metricType, metricName).Return("counter:requests").Once()
	// Mocking Set to succeed
	mockStorage.On("Set", mock.Anything, "counter:requests", metricValue).Return(nil).Once()

	// Act
	err := repo.Save(context.Background(), metricType, metricName, metricValue)

	// Assert
	assert.Nil(t, err)
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

func TestMetricRepository_Save_Failure(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	repo := &MetricRepository{StorageEngine: mockStorage, KeyEngine: mockKeyEngine}

	metricType := "counter"
	metricName := "requests"
	metricValue := "100"

	// Mocking Encode to return the correct key
	mockKeyEngine.On("Encode", metricType, metricName).Return("counter:requests").Once()
	// Mocking Set to return an error
	mockStorage.On("Set", mock.Anything, "counter:requests", metricValue).Return(errors.New("storage error")).Once()

	// Act
	err := repo.Save(context.Background(), metricType, metricName, metricValue)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "storage error", err.Error())
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

func TestMetricRepository_Get_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	repo := &MetricRepository{StorageEngine: mockStorage, KeyEngine: mockKeyEngine}

	metricType := "counter"
	metricName := "requests"
	metricValue := "100"

	// Mocking Encode to return the correct key
	mockKeyEngine.On("Encode", metricType, metricName).Return("counter:requests").Once()
	// Mocking Get to return the value successfully
	mockStorage.On("Get", mock.Anything, "counter:requests").Return(metricValue, nil).Once()

	// Act
	result, err := repo.Get(context.Background(), metricType, metricName)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, metricValue, result)
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

func TestMetricRepository_Get_Failure(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	repo := &MetricRepository{StorageEngine: mockStorage, KeyEngine: mockKeyEngine}

	metricType := "counter"
	metricName := "requests"

	// Mocking Encode to return the correct key
	mockKeyEngine.On("Encode", metricType, metricName).Return("counter:requests").Once()
	// Mocking Get to return an error
	mockStorage.On("Get", mock.Anything, "counter:requests").Return("", errors.New("storage error")).Once()

	// Act
	result, err := repo.Get(context.Background(), metricType, metricName)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, types.EmptyString, result)
	assert.Equal(t, "storage error", err.Error())
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

func TestMetricRepository_GetAll_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	repo := &MetricRepository{StorageEngine: mockStorage, KeyEngine: mockKeyEngine}

	// Preparing mock data for Generate
	ch := make(chan []string, 2)
	mockStorage.On("Generate", mock.Anything).Return(ch)

	// Simulating valid data for the test
	go func() {
		ch <- []string{"counter:requests", "100"}
		ch <- []string{"gauge:cpu_usage", "0.85"}
		close(ch)
	}()

	// Mocking Decode to return the decoded values
	mockKeyEngine.On("Decode", "counter:requests").Return("counter", "requests", nil).Once()
	mockKeyEngine.On("Decode", "gauge:cpu_usage").Return("gauge", "cpu_usage", nil).Once()

	// Act
	result := repo.GetAll(context.Background())

	// Assert
	expected := [][]string{
		{"counter", "requests", "100"},
		{"gauge", "cpu_usage", "0.85"},
	}
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

func TestMetricRepository_GetAll_DecodeError(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	repo := &MetricRepository{StorageEngine: mockStorage, KeyEngine: mockKeyEngine}

	// Preparing mock data for Generate
	ch := make(chan []string, 2)
	mockStorage.On("Generate", mock.Anything).Return(ch)

	// Simulating data with an invalid key for decoding
	go func() {
		ch <- []string{"counter:requests", "100"}
		ch <- []string{"invalid_key", "999"} // This should fail decoding
		close(ch)
	}()

	// Mocking Decode to return an error for the invalid key
	mockKeyEngine.On("Decode", "counter:requests").Return("counter", "requests", nil).Once()
	mockKeyEngine.On("Decode", "invalid_key").Return("", "", errors.New("decode error")).Once()

	// Act
	result := repo.GetAll(context.Background())

	// Assert
	expected := [][]string{
		{"counter", "requests", "100"},
	}
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}
