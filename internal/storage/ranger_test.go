package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тесты для Ranger
func TestRanger_Range(t *testing.T) {
	// Создаём хранилище и заполняем его данными
	storage := &Storage[int, string]{data: make(map[int]string)}
	storage.data[1] = "value1"
	storage.data[2] = "value2"
	storage.data[3] = "value3"

	// Создаём Ranger
	ranger := NewRanger(storage)

	// Проверяем, что callback вызывается для каждого элемента
	var keys []int
	var values []string

	ranger.Range(func(key int, value string) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true // продолжаем перебор
	})

	// Проверяем порядок элементов
	assert.Equal(t, []int{1, 2, 3}, keys, "Ключи должны быть правильно перебраны")
	assert.Equal(t, []string{"value1", "value2", "value3"}, values, "Значения должны быть правильно перебраны")

}

func TestRanger_Range_StopOnCallbackFalse(t *testing.T) {
	// Создаём хранилище и заполняем его данными
	storage := &Storage[int, string]{data: make(map[int]string)}
	storage.data[1] = "value1"
	storage.data[2] = "value2"
	storage.data[3] = "value3"

	// Создаём Ranger
	ranger := NewRanger(storage)

	// Проверяем, что перебор останавливается, когда callback возвращает false
	var keys []int
	var values []string

	ranger.Range(func(key int, value string) bool {
		keys = append(keys, key)
		values = append(values, value)
		// Прерываем перебор, когда ключ == 2
		return key != 2
	})

	// Проверяем, что перебор остановился на ключе 2
	// В результате должны быть только ключи 1 и 2, а также их значения
	assert.Equal(t, []int{1, 2}, keys, "Ключи должны быть 1 и 2")
	assert.Equal(t, []string{"value1", "value2"}, values, "Значения должны быть 'value1' и 'value2'")

	// Теперь проверим, что ключ 3 не был добавлен в результат
	assert.NotContains(t, keys, 3, "Ключ 3 не должен быть в результатах")
	assert.NotContains(t, values, "value3", "Значение 'value3' не должно быть в результатах")
}
