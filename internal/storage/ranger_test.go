package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRanger_Range(t *testing.T) {
	// Создаем хранилище и Saver для добавления данных
	storage := &Storage[int]{data: make(map[string]int)}
	saver := NewSaver(storage)

	// Очистка данных перед тестами
	t.Cleanup(func() {
		storage.data = make(map[string]int)
	})

	// Добавляем данные в хранилище
	saver.Save("key1", 1)
	saver.Save("key2", 2)
	saver.Save("key3", 3)

	// Создаем Ranger
	ranger := NewRanger(storage)

	// Тест 1: Проверяем, что Range правильно перебирает все элементы
	var result []string
	ranger.Range(func(key string, value int) bool {
		result = append(result, key)
		return true
	})

	// Проверяем, что все ключи были перебраны
	assert.ElementsMatch(t, result, []string{"key1", "key2", "key3"})

	// Тест 2: Проверяем, что Range прекращает перебор при возвращении false из callback
	result = []string{}
	ranger.Range(func(key string, value int) bool {
		if key == "key2" {
			return false
		}
		result = append(result, key)
		return true
	})

	// Проверяем, что перебор остановился на key2
	assert.ElementsMatch(t, result, []string{"key1"})
}

func TestRanger_EmptyStorage(t *testing.T) {
	// Создаем пустое хранилище
	storage := &Storage[int]{data: make(map[string]int)}
	ranger := NewRanger(storage)

	// Тест: Проверяем, что Range ничего не возвращает для пустого хранилища
	result := []string{}
	ranger.Range(func(key string, value int) bool {
		result = append(result, key)
		return true
	})

	// Хранилище пустое, значит, result должен быть пустым
	assert.Empty(t, result)
}

func TestRanger_WithOneElement(t *testing.T) {
	// Создаем хранилище с одним элементом
	storage := &Storage[int]{data: make(map[string]int)}
	saver := NewSaver(storage)
	saver.Save("key1", 1)

	// Создаем Ranger
	ranger := NewRanger(storage)

	// Тест: Проверяем, что один элемент был правильно перебран
	result := []string{}
	ranger.Range(func(key string, value int) bool {
		result = append(result, key)
		return true
	})

	// Ожидаем, что result будет содержать только один ключ "key1"
	assert.ElementsMatch(t, result, []string{"key1"})
}

func TestRanger_ConcurrentAccess(t *testing.T) {
	// Создаем хранилище с несколькими элементами
	storage := &Storage[int]{data: make(map[string]int)}
	saver := NewSaver(storage)
	saver.Save("key1", 1)
	saver.Save("key2", 2)

	// Создаем Ranger
	ranger := NewRanger(storage)

	// Тест: Проверяем, что можно безопасно работать с Ranger в конкурентной среде
	var result []string
	done := make(chan bool)
	go func() {
		ranger.Range(func(key string, value int) bool {
			result = append(result, key)
			return true
		})
		done <- true
	}()

	// Добавляем новый элемент в хранилище во время перебора
	saver.Save("key3", 3)

	// Ждем завершения перебора
	<-done

	// Проверяем, что все ключи были перебраны
	assert.ElementsMatch(t, result, []string{"key1", "key2", "key3"})
}
