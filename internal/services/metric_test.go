package services_test

import (
	"go-metrics-alerting/internal/api/types"
	"go-metrics-alerting/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGaugeSaver struct {
	mock.Mock
}

func (m *MockGaugeSaver) Save(key string, value float64) bool {
	args := m.Called(key, value)
	return args.Bool(0)
}

type MockGaugeGetter struct {
	mock.Mock
}

func (m *MockGaugeGetter) Get(key string) (float64, bool) {
	args := m.Called(key)
	return args.Get(0).(float64), args.Bool(1)
}

type MockCounterSaver struct {
	mock.Mock
}

func (m *MockCounterSaver) Save(key string, value int64) bool {
	args := m.Called(key, value)
	return args.Bool(0)
}

type MockCounterGetter struct {
	mock.Mock
}

func (m *MockCounterGetter) Get(key string) (int64, bool) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Bool(1)
}

func TestUpdateMetricsService_Update_Gauge(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := services.NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: types.Gauge,
		Value: float64Ptr(5.5),
	}

	// Expect Save to be called once for the gauge
	mockGaugeSaver.On("Save", "metric1", 5.5).Return(true)

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	assert.Equal(t, float64(5.5), *resp.Value)
	mockGaugeSaver.AssertExpectations(t)
}

func TestUpdateMetricsService_Update_Counter_ExistingValue(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := services.NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: types.Counter,
		Delta: int64Ptr(3),
	}

	// Setup mock: counter exists with value 2
	mockCounterGetter.On("Get", "metric1").Return(int64(2), true)
	// Expect Save to be called with the updated value: 2 + 3 = 5
	mockCounterSaver.On("Save", "metric1", int64(5)).Return(true)

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	mockCounterGetter.AssertExpectations(t)
	mockCounterSaver.AssertExpectations(t)
}

func TestUpdateMetricsService_Update_Counter_NoExistingValue(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := services.NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: types.Counter,
		Delta: int64Ptr(3),
	}

	// Setup mock: counter does not exist, so return default value 0
	mockCounterGetter.On("Get", "metric1").Return(int64(0), false)
	// Expect Save to be called with the updated value: 0 + 3 = 3
	mockCounterSaver.On("Save", "metric1", int64(3)).Return(true)

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	mockCounterGetter.AssertExpectations(t)
	mockCounterSaver.AssertExpectations(t)
}

func TestUpdateMetricsService_Update_UnsupportedType(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := services.NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: "Unsupported", // This is an unsupported type
	}

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, resp) // No response for unsupported types
}

func float64Ptr(v float64) *float64 {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}
