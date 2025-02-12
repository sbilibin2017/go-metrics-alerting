package engines

import (
	"context"
	"errors"
	"sync"
)

// Setter определяет интерфейс для сохранения данных.
type Setter interface {
	Set(ctx context.Context, key string, value string) error
}

// Getter определяет интерфейс для получения данных.
type Getter interface {
	Get(ctx context.Context, key string) (string, error)
}

// Generator определяет интерфейс для генерации всех данных.
type Generator interface {
	Generate(ctx context.Context) <-chan []string
}

// Storage объединяет Setter, Getter и Generator.
type Storage interface {
	Setter
	Getter
	Generator
}

const (
	StorageEmptyString string = ""
)

var (
	ErrContextDone   error = errors.New("context done")
	ErrValueNotFound error = errors.New("value not found")
)

// MemStorage реализует интерфейс Storage, используя sync.Map.
type StorageEngine struct {
	data sync.Map
}

// Set сохраняет метрику в потокобезопасном режиме.
func (m *StorageEngine) Set(ctx context.Context, key string, value string) error {
	select {
	case <-ctx.Done(): // Проверяем, отменен ли контекст.
		return ErrContextDone // Используем константу для ошибки
	default:
		m.data.Store(key, value)
		return nil
	}
}

// Get получает метрику по ключу в потокобезопасном режиме.
func (m *StorageEngine) Get(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done(): // Проверяем, отменен ли контекст.
		return StorageEmptyString, ErrContextDone // Используем константу для ошибки
	default:
		if value, ok := m.data.Load(key); ok {
			return value.(string), nil
		}
		return StorageEmptyString, ErrValueNotFound // Используем константу для ошибки
	}
}

func (m *StorageEngine) Generate(ctx context.Context) <-chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)

		m.data.Range(func(key, value any) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- []string{key.(string), value.(string)}:
				return true
			}
		})
	}()
	return ch
}
