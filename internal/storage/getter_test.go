package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тесты для Getter
func TestGetter_Get(t *testing.T) {
	// Создаём хранилище с данными
	storage := &Storage[int, string]{data: make(map[int]string)}
	storage.data[1] = "test_value"
	storage.data[2] = "another_value"

	// Создаём Getter
	getter := NewGetter(storage)

	// Проверяем получение существующего значения
	value, exists := getter.Get(1)
	assert.True(t, exists, "Значение должно существовать для ключа 1")
	assert.Equal(t, "test_value", value, "Значение для ключа 1 должно быть 'test_value'")

	// Проверяем получение несуществующего значения
	value, exists = getter.Get(3)
	assert.False(t, exists, "Значение не должно существовать для ключа 3")
	assert.Equal(t, "", value, "Значение для ключа 3 должно быть пустым")
}

func TestGetter_Concurrency(t *testing.T) {
	// Создаём хранилище и заполняем его данными
	storage := &Storage[int, string]{data: make(map[int]string)}
	storage.data[1] = "test_value"
	storage.data[2] = "another_value"

	// Создаём Getter
	getter := NewGetter(storage)

	var wg sync.WaitGroup

	// Запускаем несколько горутин для параллельного чтения из хранилища
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value, exists := getter.Get(i)
			if i == 1 {
				assert.True(t, exists, "Значение должно существовать для ключа 1")
				assert.Equal(t, "test_value", value, "Значение для ключа 1 должно быть 'test_value'")
			} else if i == 2 {
				assert.True(t, exists, "Значение должно существовать для ключа 2")
				assert.Equal(t, "another_value", value, "Значение для ключа 2 должно быть 'another_value'")
			}
		}(i)
	}

	wg.Wait()
}
