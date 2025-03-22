package registries

import (
	"context"
	"fmt"
)

// Интерфейс Worker, который реализуют воркеры
type Worker interface {
	Start(ctx context.Context) error
}

// Структура для хранения зарегистрированных воркеров
type WorkerRegistry struct {
	workers []Worker // Слайс для хранения воркеров
}

// NewWorkerRegistry создает и возвращает новый экземпляр WorkerRegistry.
func NewWorkerRegistry() *WorkerRegistry {
	return &WorkerRegistry{
		workers: make([]Worker, 0), // Инициализируем пустой слайс для хранения воркеров
	}
}

// Метод для добавления воркера в реестр
func (r *WorkerRegistry) Register(worker Worker) error {
	r.workers = append(r.workers, worker)
	return nil
}

// Метод для запуска всех зарегистрированных воркеров
func (r *WorkerRegistry) StartAll(ctx context.Context) error {
	for _, worker := range r.workers {
		if err := worker.Start(ctx); err != nil {
			return fmt.Errorf("error starting worker: %v", err)
		}
	}
	return nil
}
