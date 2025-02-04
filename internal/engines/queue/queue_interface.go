package queue

// QueueEngineInterface определяет методы для работы с очередью.
type QueueEngineInterface[T any] interface {
	Produce(item T)
	Consume() T
	Close()
}
