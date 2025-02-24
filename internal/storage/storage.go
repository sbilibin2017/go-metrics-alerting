package storage

// Storage является обобщённым хранилищем данных с синхронизацией.
type Storage[T any] struct {
	data map[string]T
}

// NewStorage создаёт и возвращает новое хранилище для данных.
func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		data: make(map[string]T),
	}
}
