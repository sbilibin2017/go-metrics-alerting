package storage

import (
	"go-metrics-alerting/internal/types"
	"sync"
)

// Storage является основным хранилищем данных с синхронизацией.
type Storage struct {
	data map[string]*types.Metrics
	mu   sync.RWMutex
}

// NewStorage создаёт и возвращает новое хранилище для данных.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]*types.Metrics),
	}
}

// Saver управляет операцией записи в хранилище.
type Saver struct {
	storage *Storage
}

func NewSaver(storage *Storage) *Saver {
	return &Saver{storage: storage}
}

func (s *Saver) Save(key string, value *types.Metrics) bool {
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

func (g *Getter) Get(key string) *types.Metrics {
	g.storage.mu.RLock()
	defer g.storage.mu.RUnlock()
	return g.storage.data[key]
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger struct {
	storage *Storage
}

func NewRanger(storage *Storage) *Ranger {
	return &Ranger{storage: storage}
}

func (r *Ranger) Range(callback func(key string, value *types.Metrics) bool) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
