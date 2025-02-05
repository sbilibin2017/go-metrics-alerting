package storage

import (
	"errors"
	"sync"
)

var ErrValueNotFound = errors.New("value not found")

// MemStorage - структура для хранения метрик в памяти
type memStorage struct {
	data map[string]string
	mu   sync.RWMutex // для защиты данных при многопоточном доступе
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *memStorage {
	return &memStorage{
		data: make(map[string]string),
	}
}

// Set сохраняет метрику в хранилище с блокировкой для обеспечения потокобезопасности
func (m *memStorage) Set(key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	return nil
}

// Get извлекает метрику по ключу с блокировкой для обеспечения потокобезопасности
func (m *memStorage) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.data[key]
	if !exists {
		return "", ErrValueNotFound
	}
	return value, nil
}

// Generate возвращает канал, из которого по очереди можно извлекать все метрики
func (m *memStorage) Generate() <-chan [2]string {
	ch := make(chan [2]string)

	go func() {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for key, value := range m.data {
			ch <- [2]string{key, value}
		}

		close(ch)
	}()

	return ch
}
