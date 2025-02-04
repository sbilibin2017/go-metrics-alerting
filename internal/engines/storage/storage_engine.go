package storage

import (
	"sync"
)

// StorageEngine — это универсальный механизм хранения данных в памяти.
type StorageEngine[K comparable, V any] struct {
	data sync.Map
}

// NewStorageEngine создает новый экземпляр StorageEngine.
func NewStorageEngine[K comparable, V any]() *StorageEngine[K, V] {
	return &StorageEngine[K, V]{}
}

// Set сохраняет значение по указанному ключу.
func (s *StorageEngine[K, V]) Set(key K, value V) {
	s.data.Store(key, value)
}

// Get получает значение по указанному ключу.
func (s *StorageEngine[K, V]) Get(key K) (V, bool) {
	val, ok := s.data.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return val.(V), true
}

// Generate возвращает канал, по которому можно итерироваться по всем ключам и значениям.
func (s *StorageEngine[K, V]) Generate() <-chan [2]V {
	ch := make(chan [2]V)

	go func() {
		s.data.Range(func(key, value any) bool {
			ch <- [2]V{value.(V)}
			return true
		})
		close(ch)
	}()

	return ch
}

// Проверка соответствия интерфейсу
var _ StorageEngineInterface[string, []byte] = (*StorageEngine[string, []byte])(nil)
