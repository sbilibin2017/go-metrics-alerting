package storage

import (
	"errors"
	"sync"
)

const (
	StorageEmptyString string = ""
)

var (
	ErrContextDone   = errors.New("context done")
	ErrValueNotFound = errors.New("value not found")
)

// Storage реализует интерфейс хранилища, используя map и sync.RWMutex.
type MemStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]string)}
}

// Set добавляет пару ключ-значение в хранилище.
func (s *MemStorage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get извлекает значение по ключу.
func (s *MemStorage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Range перебирает все пары ключ-значение.
func (s *MemStorage) Range(callback func(key, value string) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for key, value := range s.data {
		if !callback(key, value) {
			break
		}
	}
}
