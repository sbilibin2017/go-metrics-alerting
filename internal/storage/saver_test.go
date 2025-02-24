package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaver_Save(t *testing.T) {
	// Создаем хранилище для теста
	storage := NewStorage[int, string]()
	saver := NewSaver(storage)

	// Тестируем сохранение данных
	saver.Save(1, "test_value")
	assert.Equal(t, "test_value", storage.data[1], "Значение должно быть сохранено в хранилище")

	// Тестируем замену значения по ключу
	saver.Save(1, "new_value")
	assert.Equal(t, "new_value", storage.data[1], "Значение должно быть обновлено в хранилище")

	// Проверка сохранения другого ключа
	saver.Save(2, "another_value")
	assert.Equal(t, "another_value", storage.data[2], "Значение для другого ключа должно быть сохранено корректно")
}

func TestSaver_Concurrency(t *testing.T) {
	storage := NewStorage[int, string]()
	saver := NewSaver(storage)

	// Используем несколько горутин для проверки корректности работы в многозадачной среде
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			saver.Save(i, "value")
		}(i)
	}

	wg.Wait()

	// Проверка, что все значения были сохранены
	for i := 0; i < 1000; i++ {
		assert.Equal(t, "value", storage.data[i], "Значение для ключа %d должно быть равно 'value'", i)
	}
}
