package storage

import "sync"

// Saver управляет операцией записи в обобщённое хранилище.
type Saver[T any] struct {
	storage *Storage[T]
	mu      sync.RWMutex
}

// NewSaver создаёт новый экземпляр Saver для работы с хранилищем типа T.
func NewSaver[T any](storage *Storage[T]) *Saver[T] {
	return &Saver[T]{storage: storage}
}

// Save сохраняет данные в хранилище по указанному ключу.
func (s *Saver[T]) Save(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage.data[key] = value
}
