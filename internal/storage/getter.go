package storage

import "sync"

// Getter управляет операцией чтения из обобщённого хранилища с двумя параметрами типа K и V.
type Getter[K comparable, V any] struct {
	storage *Storage[K, V]
	mu      sync.RWMutex
}

// NewGetter создаёт новый экземпляр Getter для работы с хранилищем типа K и V.
func NewGetter[K comparable, V any](storage *Storage[K, V]) *Getter[K, V] {
	return &Getter[K, V]{storage: storage}
}

// Get получает данные из хранилища по указанному ключу.
func (g *Getter[K, V]) Get(key K) (V, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	value, exists := g.storage.data[key]
	return value, exists
}
