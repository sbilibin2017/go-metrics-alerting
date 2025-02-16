package storage

import (
	"errors"
	"sync"
)

var ErrNotFound error = errors.New("not found")

const EmptyString string = ""

// Storage является основным хранилищем данных с синхронизацией для строковых значений.
type Storage struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewStorage создаёт и возвращает новое хранилище для строковых данных.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Setter управляет операцией записи в хранилище.
type Setter struct {
	storage *Storage
}

func NewSetter(storage *Storage) *Setter {
	return &Setter{storage: storage}
}

func (s *Setter) Set(key string, value string) error {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.data[key] = value
	return nil
}

// Getter управляет операцией чтения из хранилища.
type Getter struct {
	storage *Storage
}

func NewGetter(storage *Storage) *Getter {
	return &Getter{storage: storage}
}

func (g *Getter) Get(key string) (string, error) {
	g.storage.mu.RLock()
	defer g.storage.mu.RUnlock()
	value, exists := g.storage.data[key]
	if !exists {
		return EmptyString, ErrNotFound
	}
	return value, nil
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger struct {
	storage *Storage
}

func NewRanger(storage *Storage) *Ranger {
	return &Ranger{storage: storage}
}

func (r *Ranger) Range(callback func(key string, value string) bool) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
