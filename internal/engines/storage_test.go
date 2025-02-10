package engines

import (
	"context"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_Set(t *testing.T) {
	storage := &StorageEngine{}

	tests := []struct {
		key   string
		value string
	}{
		{"metric1", "10"},
		{"metric2", "20"},
		{"metric3", "30"},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			// Проверяем, что ошибка нет
			err := storage.Set(context.Background(), test.key, test.value)
			assert.NoError(t, err)
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	storage := &StorageEngine{}
	storage.Set(context.Background(), "metric1", "10")
	storage.Set(context.Background(), "metric2", "20")

	tests := []struct {
		key         string
		expected    string
		expectedErr error
	}{
		{"metric1", "10", nil},
		{"metric2", "20", nil},
		{"metric3", types.EmptyString, errors.ErrValueNotFound}, // Используем EmptyString и ErrValueNotFound
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			result, err := storage.Get(context.Background(), test.key)

			// Проверяем, что ошибка соответствует ожидаемой
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем, что возвращаемое значение совпадает с ожидаемым
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMemStorage_Generate(t *testing.T) {
	storage := &StorageEngine{}
	storage.Set(context.Background(), "metric1", "10")
	storage.Set(context.Background(), "metric2", "20")

	tests := []struct {
		expectedKeys   []string
		expectedValues []string
	}{
		{[]string{"metric1", "metric2"}, []string{"10", "20"}},
	}

	for _, test := range tests {
		t.Run("Generate", func(t *testing.T) {
			ch := storage.Generate(context.Background())

			// Временно собираем результаты из канала
			var keys []string
			var values []string

			for kv := range ch {
				keys = append(keys, kv[0])
				values = append(values, kv[1])
			}

			// Проверяем, что ключи и значения соответствуют ожидаемым
			assert.ElementsMatch(t, test.expectedKeys, keys)
			assert.ElementsMatch(t, test.expectedValues, values)
		})
	}
}

func TestMemStorage_Set_ContextCanceled(t *testing.T) {
	storage := &StorageEngine{}
	ctx, cancel := context.WithCancel(context.Background())

	// Cancelling the context to simulate an error
	cancel()

	err := storage.Set(ctx, "metric1", "10")

	// Проверяем, что ошибка соответствует ErrContextDone
	assert.EqualError(t, err, errors.ErrContextDone.Error())
}

func TestMemStorage_Get_ContextCanceled(t *testing.T) {
	storage := &StorageEngine{}
	storage.Set(context.Background(), "metric1", "10")
	ctx, cancel := context.WithCancel(context.Background())

	// Cancelling the context to simulate an error
	cancel()

	_, err := storage.Get(ctx, "metric1")

	// Проверяем, что ошибка соответствует ErrContextDone
	assert.EqualError(t, err, errors.ErrContextDone.Error())
}

func TestMemStorage_Generate_ContextCanceled(t *testing.T) {
	storage := &StorageEngine{}
	storage.Set(context.Background(), "metric1", "10")
	storage.Set(context.Background(), "metric2", "20")
	ctx, cancel := context.WithCancel(context.Background())

	// Cancelling the context to stop the generator early
	cancel()

	ch := storage.Generate(ctx)

	// Проверяем, что канал закрывается сразу после отмены контекста
	select {
	case _, ok := <-ch:
		assert.False(t, ok, "expected channel to be closed due to canceled context")
	default:
		// No output expected as the context is canceled
	}
}
