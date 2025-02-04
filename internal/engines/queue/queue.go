package queue

// QueueEngine реализует очередь с использованием канала.
type QueueEngine[T any] struct {
	queue chan T
}

// NewQueueEngine создает новый экземпляр очереди с заданной емкостью.
func NewQueueEngine[T any](capacity int) *QueueEngine[T] {
	return &QueueEngine[T]{
		queue: make(chan T, capacity),
	}
}

// Produce добавляет элемент в очередь.
func (qe *QueueEngine[T]) Produce(item T) {
	qe.queue <- item
}

// Consume извлекает элемент из очереди.
func (qe *QueueEngine[T]) Consume() T {
	item := <-qe.queue
	return item
}

// Close закрывает канал очереди.
func (qe *QueueEngine[T]) Close() {
	close(qe.queue)
}

// Проверка соответствия интерфейсу
var _ QueueEngineInterface[any] = (*QueueEngine[any])(nil)
