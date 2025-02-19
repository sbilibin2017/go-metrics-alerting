package services

import (
	"go-metrics-alerting/internal/server/types"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(key, value string) bool {
	args := m.Called(key, value)
	return args.Bool(0)
}

type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

type MockRanger struct {
	mock.Mock
}

func (m *MockRanger) Range(callback func(key, value string) bool) {
	args := m.Called()
	data := args.Get(0).(map[string]string)
	for k, v := range data {
		if !callback(k, v) {
			break
		}
	}
}

func TestUpdateMetricsService_Update_Gauge(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	service := NewUpdateMetricsService(saver, getter)

	// Test Gauge
	value := 42.5
	reqGauge := &types.UpdateMetricsRequest{ID: "test_gauge", MType: types.Gauge, Value: &value}
	saver.On("Save", reqGauge.ID, strconv.FormatFloat(value, 'f', -1, 64)).Return(true)

	// Test only the first assertion
	_, err := service.Update(reqGauge)

	assert.NoError(t, err)
}

func TestUpdateMetricsService_Update_Counter(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	service := NewUpdateMetricsService(saver, getter)

	// Test Counter
	initialValue := "10"
	getter.On("Get", "test_counter").Return(initialValue, true)
	counterDelta := int64(5)
	reqCounter := &types.UpdateMetricsRequest{ID: "test_counter", MType: types.Counter, Delta: &counterDelta}
	saver.On("Save", reqCounter.ID, "15").Return(true)

	// Test only the first assertion
	_, err := service.Update(reqCounter)

	assert.NoError(t, err)
}

func TestUpdateMetricsService_Update_Invalid(t *testing.T) {
	saver := new(MockSaver)
	getter := new(MockGetter)
	service := NewUpdateMetricsService(saver, getter)

	// Test Default Case
	reqInvalid := &types.UpdateMetricsRequest{ID: "invalid_metric", MType: "invalid_type"}
	resp, _ := service.Update(reqInvalid)

	// Test only the first assertion
	assert.Nil(t, resp)
}

func TestGetMetricValueService_GetMetricValue_Valid(t *testing.T) {
	getter := new(MockGetter)
	service := NewGetMetricValueService(getter)

	getter.On("Get", "test_metric").Return("100", true)

	req := &types.GetMetricValueRequest{ID: "test_metric"}
	_, err := service.GetMetricValue(req)

	// Test only the first assertion
	assert.NoError(t, err)
}

func TestGetMetricValueService_GetMetricValue_Invalid(t *testing.T) {
	getter := new(MockGetter)
	service := NewGetMetricValueService(getter)

	// Test Default Case
	reqInvalid := &types.GetMetricValueRequest{ID: "invalid_metric"}
	getter.On("Get", reqInvalid.ID).Return("", false)
	resp, _ := service.GetMetricValue(reqInvalid)

	// Test only the first assertion
	assert.Nil(t, resp)
}

func TestGetAllMetricValuesService_GetAllMetricValues(t *testing.T) {
	ranger := new(MockRanger)
	service := NewGetAllMetricValuesService(ranger)

	metricsData := map[string]string{
		"metric1": "100",
		"metric2": "200",
	}
	ranger.On("Range").Return(metricsData)

	// Test only the first assertion
	resp := service.GetAllMetricValues()

	assert.Len(t, resp, 2)
}
