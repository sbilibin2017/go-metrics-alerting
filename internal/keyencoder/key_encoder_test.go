package keyencoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEncoder_Encode(t *testing.T) {
	// Создание объекта для тестирования
	encoder := NewKeyEncoder()

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
