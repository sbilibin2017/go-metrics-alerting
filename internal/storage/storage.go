package storage

import (
	"sync"
)

// Storage является основным хранилищем данных с синхронизацией для любых типов значений.
type Storage[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

// NewStorage создаёт и возвращает новое хранилище для данных.
func NewStorage[K comparable, V any]() *Storage[K, V] {
	return &Storage[K, V]{
		data: make(map[K]V),
	}
}

// Saver управляет операцией записи в хранилище.
type Saver[K comparable, V any] struct {
	storage *Storage[K, V]
}

func NewSaver[K comparable, V any](storage *Storage[K, V]) *Saver[K, V] {
	return &Saver[K, V]{storage: storage}
}

func (s *Saver[K, V]) Save(key K, value V) bool {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.data[key] = value
	return true
}

// Getter управляет операцией чтения из хранилища.
type Getter[K comparable, V any] struct {
	storage *Storage[K, V]
}

func NewGetter[K comparable, V any](storage *Storage[K, V]) *Getter[K, V] {
	return &Getter[K, V]{storage: storage}
}

func (g *Getter[K, V]) Get(key K) (V, bool) {
	g.storage.mu.RLock()
	defer g.storage.mu.RUnlock()
	value, exists := g.storage.data[key]
	var zeroValue V
	if !exists {
		return zeroValue, false
	}
	return value, true
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger[K comparable, V any] struct {
	storage *Storage[K, V]
}

func NewRanger[K comparable, V any](storage *Storage[K, V]) *Ranger[K, V] {
	return &Ranger[K, V]{storage: storage}
}

func (r *Ranger[K, V]) Range(callback func(key K, value V) bool) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
