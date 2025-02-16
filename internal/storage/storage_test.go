package storage

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Setter_Set(t *testing.T) {
	// Создаем хранилище и репозиторий для записи
	storage := NewStorage()
	setter := NewSetter(storage)

	// Выполняем операцию записи
	setter.Set("metricType1:metricName1", "metricValue1")

	// Проверяем, что данные были записаны в хранилище
	getter := NewGetter(storage)
	value, exists := getter.Get("metricType1:metricName1")
	assert.True(t, exists, "Expected value to exist")
	assert.Equal(t, "metricValue1", value, "Expected value to be 'metricValue1'")
}

func TestStorage_Setter_Set_Concurrency(t *testing.T) {
	// Создаем хранилище и репозиторий для записи
	storage := NewStorage()
	setter := NewSetter(storage)

	// Используем goroutines для конкурентной записи
	done := make(chan struct{})
	for i := 0; i < 1000; i++ {
		go func(i int) {
			// Исправлено: используем strconv.Itoa для преобразования числа в строку
			setter.Set("key", "value"+strconv.Itoa(i))
			done <- struct{}{}
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 1000; i++ {
		<-done
	}

	// Проверяем, что после конкурентной записи значение в хранилище присутствует
	getter := NewGetter(storage)
	value, exists := getter.Get("key")
	assert.True(t, exists, "Expected value to exist")
	assert.Contains(t, value, "value", "Expected value to contain 'value'")
}

func TestStorage_Getter_Get(t *testing.T) {
	// Создаем хранилище и репозиторий для записи и чтения
	storage := NewStorage()
	setter := NewSetter(storage)
	getter := NewGetter(storage)

	// Записываем данные в хранилище
	setter.Set("metricType1:metricName1", "metricValue1")

	// Выполняем операцию чтения
	value, exists := getter.Get("metricType1:metricName1")

	// Проверяем, что данные существуют и корректны
	assert.True(t, exists, "Expected value to exist")
	assert.Equal(t, "metricValue1", value, "Expected value to be 'metricValue1'")

	// Пытаемся получить несуществующее значение
	_, exists = getter.Get("nonExistentKey")
	assert.False(t, exists, "Expected value to not exist")
}

func TestStorage_Ranger_Range(t *testing.T) {
	// Создаем хранилище и репозиторий для записи и перебора
	storage := NewStorage()
	setter := NewSetter(storage)
	ranger := NewRanger(storage)

	// Записываем данные в хранилище
	setter.Set("metricType1:metricName1", "metricValue1")
	setter.Set("metricType2:metricName2", "metricValue2")

	// Перебираем данные и проверяем, что они правильные
	var result []string
	ranger.Range(func(key, value string) bool {
		result = append(result, key+":"+value)
		return true
	})

	// Проверяем количество элементов и их содержание
	assert.Equal(t, 2, len(result), "Expected to iterate over 2 elements")
	assert.Contains(t, result, "metricType1:metricName1:metricValue1", "Expected key-value pair 'metricType1:metricName1' to be present")
	assert.Contains(t, result, "metricType2:metricName2:metricValue2", "Expected key-value pair 'metricType2:metricName2' to be present")
}

func TestStorage_Getter_Get_NonExistentKey(t *testing.T) {
	// Создаем хранилище и репозиторий для чтения
	storage := NewStorage()
	getter := NewGetter(storage)

	// Пытаемся получить несуществующий ключ
	value, exists := getter.Get("nonExistentKey")

	// Проверяем, что значение не найдено
	assert.False(t, exists, "Expected value to not exist")
	assert.Equal(t, "", value, "Expected value to be empty string")
}

func TestStorage_Ranger_Range_Empty(t *testing.T) {
	// Создаем пустое хранилище и репозиторий для перебора
	storage := NewStorage()
	ranger := NewRanger(storage)

	// Перебираем данные (их нет) и проверяем, что результат пустой
	var result []string
	ranger.Range(func(key, value string) bool {
		result = append(result, key+":"+value)
		return true
	})

	// Проверяем, что данных нет
	assert.Empty(t, result, "Expected no data in range")
}

func TestStorage_Ranger_Range_BreakOnCallbackFalse(t *testing.T) {
	// Создаем хранилище и репозиторий для записи и перебора
	storage := NewStorage()
	setter := NewSetter(storage)
	ranger := NewRanger(storage)

	// Записываем несколько данных в хранилище
	setter.Set("metricType1:metricName1", "metricValue1")
	setter.Set("metricType2:metricName2", "metricValue2")
	setter.Set("metricType3:metricName3", "metricValue3")

	// Переменная для проверки сколько элементов было обработано
	var result []string

	// Колбэк, который останавливает перебор после 1-го элемента
	callback := func(key, value string) bool {
		result = append(result, key+":"+value)
		// Возвращаем false после первого элемента, чтобы остановить перебор
		return len(result) < 1
	}

	// Запускаем метод Range с колбэком
	ranger.Range(callback)

	// Проверяем, что в результате только один элемент
	assert.Equal(t, 1, len(result), "Expected to process only 1 element")
	assert.Contains(t, result, "metricType1:metricName1:metricValue1", "Expected first element to be 'metricType1:metricName1'")
}
