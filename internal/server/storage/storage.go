package storage

import (
	"sync"
)

// Storage является основным хранилищем данных с синхронизацией.
type Storage struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewStorage создаёт и возвращает новое хранилище для данных.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Saver управляет операцией записи в хранилище.
type Saver struct {
	storage *Storage
}

func NewSaver(storage *Storage) *Saver {
	return &Saver{storage: storage}
}

func (s *Saver) Save(key, value string) bool {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.data[key] = value
	return true
}

// Getter управляет операцией чтения из хранилища.
type Getter struct {
	storage *Storage
}

func NewGetter(storage *Storage) *Getter {
	return &Getter{storage: storage}
}

func (g *Getter) Get(key string) (string, bool) {
	g.storage.mu.RLock()
	defer g.storage.mu.RUnlock()
	value, exists := g.storage.data[key]
	return value, exists
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger struct {
	storage *Storage
}

func NewRanger(storage *Storage) *Ranger {
	return &Ranger{storage: storage}
}

func (r *Ranger) Range(callback func(key, value string) bool) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
