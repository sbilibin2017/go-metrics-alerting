package storage

import (
	"sync"
)

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger[T any] struct {
	storage *Storage[T]
	mu      sync.RWMutex
}

// NewRanger создаёт новый экземпляр Ranger для работы с хранилищем типа T.
func NewRanger[T any](storage *Storage[T]) *Ranger[T] {
	return &Ranger[T]{storage: storage}
}

// Range перебирает все элементы в хранилище и вызывает callback для каждого из них.
func (r *Ranger[T]) Range(callback func(key string, value T) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
