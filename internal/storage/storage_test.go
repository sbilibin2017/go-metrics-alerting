package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewMemStorage(t *testing.T) {
	storage := NewMemStorage()

	if storage == nil {
		t.Fatal("expected new storage, got nil")
	}
}

func TestSetAndGet(t *testing.T) {
	storage := NewMemStorage()

	// Сначала устанавливаем метрику
	key := "metric1"
	value := "100"
	if err := storage.Set(key, value); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем, что можно получить метрику
	got, err := storage.Get(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != value {
		t.Fatalf("expected %v, got %v", value, got)
	}
}

func TestGetValueNotFound(t *testing.T) {
	storage := NewMemStorage()

	// Проверяем, что возвращается ошибка, если метрика не найдена
	_, err := storage.Get("non_existing_key")
	if err != ErrValueNotFound {
		t.Fatalf("expected %v, got %v", ErrValueNotFound, err)
	}
}

func TestGenerate(t *testing.T) {
	storage := NewMemStorage()

	// Добавляем несколько метрик
	storage.Set("metric1", "100")
	storage.Set("metric2", "200")

	// Генерируем метрики
	ch := storage.Generate()

	// Проверяем, что все метрики можно получить через канал
	expected := map[string]string{
		"metric1": "100",
		"metric2": "200",
	}

	timeout := time.After(1 * time.Second)
	for len(expected) > 0 {
		select {
		case metric := <-ch:
			key, value := metric[0], metric[1]
			if expected[key] != value {
				t.Errorf("expected %v:%v, got %v:%v", key, expected[key], key, value)
			}
			delete(expected, key)
		case <-timeout:
			t.Fatal("timed out waiting for metrics")
		}
	}
}

func TestConcurrency(t *testing.T) {
	storage := NewMemStorage()

	// Записываем и извлекаем метрики одновременно
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()
			if err := storage.Set(fmt.Sprintf("%d", i), "value"); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			_, err := storage.Get(fmt.Sprintf("%d", i))
			if err != nil && err != ErrValueNotFound {
				t.Errorf("unexpected error: %v", err)
			}
		}(i)
	}

	wg.Wait()
}
