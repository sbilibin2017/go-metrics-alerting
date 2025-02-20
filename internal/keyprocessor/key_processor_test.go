package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEncoder_Encode(t *testing.T) {
	// Создание объекта для тестирования
	encoder := &KeyEncoder{}

	// Тестовые данные
	id := "metric1"
	mtype := "gauge"

	// Ожидаемый результат
	expected := "metric1:gauge"

	// Тестирование метода Encode
	result := encoder.Encode(id, mtype)

	// Сравнение результатов
	assert.Equal(t, expected, result, "Encoded key should match the expected format")
}

func TestKeyDecoder_Decode_ValidKey(t *testing.T) {
	// Создание объекта для тестирования
	decoder := &KeyDecoder{}

	// Тестовый ключ с правильным форматом
	key := "metric1:gauge"

	// Ожидаемый результат
	expectedID := "metric1"
	expectedMType := "gauge"

	// Тестирование метода Decode
	id, mtype, ok := decoder.Decode(key)

	// Проверка результата
	assert.True(t, ok, "Decode should return true for valid keys")
	assert.Equal(t, expectedID, id, "Decoded id should match the expected id")
	assert.Equal(t, expectedMType, mtype, "Decoded mtype should match the expected mtype")
}

func TestKeyDecoder_Decode_InvalidKey(t *testing.T) {
	// Создание объекта для тестирования
	decoder := &KeyDecoder{}

	// Тестовый ключ с неправильным форматом
	key := "metric1"

	// Тестирование метода Decode
	id, mtype, ok := decoder.Decode(key)

	// Проверка результата
	assert.False(t, ok, "Decode should return false for invalid keys")
	assert.Empty(t, id, "Decoded id should be empty for invalid keys")
	assert.Empty(t, mtype, "Decoded mtype should be empty for invalid keys")
}
