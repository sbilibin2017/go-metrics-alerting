package storage

import (
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSaver_Save проверяет правильность сохранения метрик в хранилище
func TestSaver_Save(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)

	// Создаём метрику напрямую
	metric := types.Metrics{
		ID:    "metric1",
		MType: types.Gauge,
		Value: new(float64),
	}
	*metric.Value = 10.5

	// Сохраняем метрику
	result := saver.Save("key1", &metric)

	// Проверяем, что результат сохранения успешен
	assert.True(t, result)

	// Проверяем, что метрика сохранена в хранилище
	storedMetric, exists := storage.data["key1"]
	assert.True(t, exists)
	assert.Equal(t, "metric1", storedMetric.ID)
	assert.Equal(t, types.Gauge, storedMetric.MType)
	assert.Equal(t, 10.5, *storedMetric.Value)
	assert.Nil(t, storedMetric.Delta)
}

// TestGetter_Get проверяет получение метрик из хранилища
func TestGetter_Get(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)
	getter := NewGetter(storage)

	// Создаём метрику напрямую
	metric := types.Metrics{
		ID:    "metric2",
		MType: types.Counter,
		Delta: new(int64),
	}
	*metric.Delta = 100

	// Сохраняем метрику
	saver.Save("key2", &metric)

	// Получаем метрику
	storedMetric := getter.Get("key2")

	// Проверяем, что метрика существует и данные правильные
	assert.NotNil(t, storedMetric)
	assert.Equal(t, "metric2", storedMetric.ID)
	assert.Equal(t, types.Counter, storedMetric.MType)
	assert.Equal(t, int64(100), *storedMetric.Delta)
	assert.Nil(t, storedMetric.Value)
}

// TestRanger_Range_WithPointer проверяет перебор данных с помощью Ranger
func TestRanger_Range_WithPointer(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)
	ranger := NewRanger(storage)

	// Сохраняем несколько данных
	metric1 := &types.Metrics{
		ID:    "metric1",
		MType: types.Gauge,
		Value: new(float64),
	}
	*metric1.Value = 10.5
	saver.Save("key1", metric1)

	metric2 := &types.Metrics{
		ID:    "metric2",
		MType: types.Counter,
		Delta: new(int64),
	}
	*metric2.Delta = 20
	saver.Save("key2", metric2)

	// Переменные для проверки, что все ключи и значения правильно обработаны
	var keys []string
	var values []*types.Metrics

	// Используем callback с логикой перебора
	ranger.Range(func(key string, value *types.Metrics) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true
	})

	// Проверяем, что все ключи и значения были правильно добавлены в список
	assert.Len(t, keys, 2)
	assert.Len(t, values, 2)

	// Проверяем порядок ключей
	assert.Equal(t, "key1", keys[0])
	assert.Equal(t, "key2", keys[1])

	// Проверяем значения метрик
	assert.Equal(t, "metric1", values[0].ID)
	assert.Equal(t, types.Gauge, values[0].MType)
	assert.Equal(t, 10.5, *values[0].Value)
	assert.Nil(t, values[0].Delta)

	assert.Equal(t, "metric2", values[1].ID)
	assert.Equal(t, types.Counter, values[1].MType)
	assert.Equal(t, int64(20), *values[1].Delta)
	assert.Nil(t, values[1].Value)
}

func TestRanger_Range_BreakCallback(t *testing.T) {
	// Создаем новое хранилище, saver и ranger
	storage := NewStorage()
	saver := NewSaver(storage)
	ranger := NewRanger(storage)

	// Сохраняем несколько метрик
	metric1 := &types.Metrics{
		ID:    "metric1",
		MType: types.Gauge,
		Value: new(float64),
	}
	*metric1.Value = 10.5
	saver.Save("key1", metric1)

	metric2 := &types.Metrics{
		ID:    "metric2",
		MType: types.Counter,
		Delta: new(int64),
	}
	*metric2.Delta = 20
	saver.Save("key2", metric2)

	// Переменные для проверки, что перебор был прерван
	var keys []string
	var values []*types.Metrics
	var calledBeforeBreak bool

	// Callback с логикой прерывания
	ranger.Range(func(key string, value *types.Metrics) bool {
		// Проверяем, что данные были переданы в callback
		keys = append(keys, key)
		values = append(values, value)

		// Логика для прерывания перебора
		if key == "key1" {
			calledBeforeBreak = true
			return false // Прерываем перебор
		}
		return true
	})

	// Проверяем, что перебор был прерван после первого элемента
	assert.True(t, calledBeforeBreak) // Убедитесь, что callback был вызван до break
	assert.Len(t, keys, 1)            // Перебор должен завершиться после первого элемента
	assert.Len(t, values, 1)          // Перебор должен завершиться после первого элемента

	// Проверяем, что добавлен только первый элемент
	assert.Equal(t, "key1", keys[0])
	assert.Equal(t, "metric1", values[0].ID)
	assert.Equal(t, types.Gauge, values[0].MType)
	assert.Equal(t, 10.5, *values[0].Value)
	assert.Nil(t, values[0].Delta)
}
