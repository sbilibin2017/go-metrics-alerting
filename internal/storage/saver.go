package storage

import "sync"

// Saver управляет операцией записи в обобщённое хранилище с двумя параметрами типа K и V.
type Saver[K comparable, V any] struct {
	storage *Storage[K, V]
	mu      sync.RWMutex
}

// NewSaver создаёт новый экземпляр Saver для работы с хранилищем типа K и V.
func NewSaver[K comparable, V any](storage *Storage[K, V]) *Saver[K, V] {
	return &Saver[K, V]{storage: storage}
}

// Save сохраняет данные в хранилище по указанному ключу.
func (s *Saver[K, V]) Save(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage.data[key] = value
}
