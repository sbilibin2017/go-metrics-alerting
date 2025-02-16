package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetter_SetAndGetter_Get(t *testing.T) {
	storage := NewStorage()
	setter := NewSetter(storage)
	getter := NewGetter(storage)

	// Устанавливаем значение
	err := setter.Set("key1", "value1")
	assert.NoError(t, err, "Ошибка при установке значения")

	// Получаем значение
	value, err := getter.Get("key1")
	assert.NoError(t, err, "Ошибка при получении существующего ключа")
	assert.Equal(t, "value1", value, "Полученное значение не совпадает с ожидаемым")

	// Получаем несуществующий ключ
	_, err = getter.Get("nonexistent")
	assert.ErrorIs(t, err, ErrNotFound, "Ожидалась ошибка ErrNotFound при запросе несуществующего ключа")
}

func TestRanger_Range(t *testing.T) {
	storage := NewStorage()
	setter := NewSetter(storage)
	ranger := NewRanger(storage)

	// Добавляем несколько значений
	_ = setter.Set("key1", "value1")
	_ = setter.Set("key2", "value2")
	_ = setter.Set("key3", "value3")

	// Проверяем перебор всех значений
	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	actual := make(map[string]string)
	ranger.Range(func(key, value string) bool {
		actual[key] = value
		return true
	})

	assert.Equal(t, expected, actual, "Перебранные значения не совпадают с ожидаемыми")
}

func TestRanger_RangeEarlyExit(t *testing.T) {
	storage := NewStorage()
	setter := NewSetter(storage)
	ranger := NewRanger(storage)

	// Добавляем значения
	_ = setter.Set("key1", "value1")
	_ = setter.Set("key2", "value2")
	_ = setter.Set("key3", "value3")

	// Проверяем, что Range выходит при false
	var count int
	ranger.Range(func(key, value string) bool {
		count++
		return count < 2 // Прерываем после первой итерации
	})

	assert.Equal(t, 2, count, "Range не завершился после указанного количества итераций")
}
