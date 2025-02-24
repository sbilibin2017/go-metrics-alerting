package storage

import (
	"sync"
)

// Ranger управляет операцией перебора элементов в хранилище с обобщёнными типами K и V.
type Ranger[K comparable, V any] struct {
	storage *Storage[K, V]
	mu      sync.RWMutex
}

// NewRanger создаёт новый экземпляр Ranger для работы с хранилищем типа K и V.
func NewRanger[K comparable, V any](storage *Storage[K, V]) *Ranger[K, V] {
	return &Ranger[K, V]{storage: storage}
}

// Range перебирает все элементы в хранилище и вызывает callback для каждого из них.
func (r *Ranger[K, V]) Range(callback func(key K, value V) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
