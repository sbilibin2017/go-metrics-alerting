package engines

import (
	"sync"
)

// StorageEngineInterface defines methods for interacting with storage
// It abstracts the storage implementation
type StorageEngineInterface interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Generate() <-chan [2]string
}

// MemoryStorageEngine implements StorageEngineInterface and stores data in memory
type MemoryStorageEngine struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemoryStorageEngine creates a new instance of MemoryStorageEngine
func NewMemoryStorageEngine() *MemoryStorageEngine {
	return &MemoryStorageEngine{
		data: make(map[string]string),
	}
}

// Set saves a value by key
func (s *MemoryStorageEngine) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get retrieves a value by key
func (s *MemoryStorageEngine) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Generate returns a channel that iterates over all key-value pairs
func (s *MemoryStorageEngine) Generate() <-chan [2]string {
	ch := make(chan [2]string)
	go func() {
		s.mu.RLock()
		defer s.mu.RUnlock()
		for key, value := range s.data {
			ch <- [2]string{key, value}
		}
		close(ch)
	}()
	return ch
}

var _ StorageEngineInterface = &MemoryStorageEngine{}
