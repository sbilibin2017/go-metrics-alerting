package storage

import (
	"sync"
)

// Getter управляет операцией чтения из обобщённого хранилища.
type Getter[T any] struct {
	storage *Storage[T]
	mu      sync.RWMutex
}

// NewGetter создаёт новый экземпляр Getter для работы с хранилищем типа T.
func NewGetter[T any](storage *Storage[T]) *Getter[T] {
	return &Getter[T]{storage: storage}
}

// Get получает данные из хранилища по указанному ключу.
func (g *Getter[T]) Get(key string) (T, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	value, exists := g.storage.data[key]
	return value, exists
}
