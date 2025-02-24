package storage

// Storage является обобщённым хранилищем данных с синхронизацией.
type Storage[K comparable, V any] struct {
	data map[K]V
}

// NewStorage создаёт и возвращает новое хранилище для данных.
func NewStorage[K comparable, V any]() *Storage[K, V] {
	return &Storage[K, V]{
		data: make(map[K]V),
	}
}
