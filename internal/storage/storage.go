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
