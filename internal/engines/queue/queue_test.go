package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueEngine_ProduceConsume(t *testing.T) {
	// Создаем очередь с емкостью 3.
	q := NewQueueEngine[int](100)
	q.Produce(1)
	q.Produce(2)
	q.Produce(3)

	// Извлекаем элементы из очереди и проверяем их.
	assert.Equal(t, 1, q.Consume())
	assert.Equal(t, 2, q.Consume())
	assert.Equal(t, 3, q.Consume())

	q.Close()
	assert.NotPanics(t, func() {
		assert.Equal(t, 0, q.Consume())
		assert.Equal(t, 0, q.Consume())
	}, "Closing the queue should not panic")
}
