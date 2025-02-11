package engines

import (
	"context"

	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageEngine_Set(t *testing.T) {
	se := &StorageEngine{data: sync.Map{}}
	ctx := context.Background()

	err := se.Set(ctx, "key1", "value1")
	assert.NoError(t, err)
}

func TestStorageEngine_Set_ContextCanceled(t *testing.T) {
	se := &StorageEngine{data: sync.Map{}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := se.Set(ctx, "key1", "value1")
	assert.Equal(t, ErrContextDone, err)
}

func TestStorageEngine_Get(t *testing.T) {
	se := &StorageEngine{data: sync.Map{}}
	ctx := context.Background()

	// Test getting a non-existing key
	value, err := se.Get(ctx, "nonexistent")
	assert.Equal(t, StorageEmptyString, value)
	assert.Equal(t, ErrValueNotFound, err)

	// Test setting and then getting a key
	se.Set(ctx, "key1", "value1")
	value, err = se.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)
}

func TestStorageEngine_Get_ContextCanceled(t *testing.T) {
	se := &StorageEngine{data: sync.Map{}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	value, err := se.Get(ctx, "key1")
	assert.Equal(t, StorageEmptyString, value)
	assert.Equal(t, ErrContextDone, err)
}

func TestStorageEngine_Generate(t *testing.T) {
	storage := &StorageEngine{}
	storage.Set(context.Background(), "cpu", "50")
	storage.Set(context.Background(), "ram", "80")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := storage.Generate(ctx)

	results := make(map[string]string)
	for metric := range ch {
		results[metric[0]] = metric[1]
	}

	assert.Equal(t, "50", results["cpu"], "CPU метрика должна быть 50")
	assert.Equal(t, "80", results["ram"], "RAM метрика должна быть 80")
}

func TestStorageEngine_Generate_ContextCanceled(t *testing.T) {
	storage := &StorageEngine{}
	storage.Set(context.Background(), "cpu", "50")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	ch := storage.Generate(ctx)

	_, ok := <-ch // Канал должен быть закрыт
	assert.False(t, ok, "Канал должен быть закрыт, так как контекст отменен")
}
