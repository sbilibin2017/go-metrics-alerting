package engines

import (
	"testing"
)

func TestMemoryStorageEngine(t *testing.T) {
	storage := NewMemoryStorageEngine()

	tests := []struct {
		name      string
		key       string
		value     string
		expectGet string
		expectOk  bool
		setBefore bool
	}{
		{
			name:      "Set and Get value",
			key:       "foo",
			value:     "bar",
			expectGet: "bar",
			expectOk:  true,
			setBefore: true,
		},
		{
			name:      "Get non-existing key",
			key:       "unknown",
			expectGet: "",
			expectOk:  false,
			setBefore: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setBefore {
				storage.Set(test.key, test.value)
			}
			result, ok := storage.Get(test.key)
			if ok != test.expectOk || result != test.expectGet {
				t.Errorf("expected (%v, %v), got (%v, %v)", test.expectGet, test.expectOk, result, ok)
			}
		})
	}
}

func TestMemoryStorageEngineGenerate(t *testing.T) {
	storage := NewMemoryStorageEngine()
	storage.Set("key1", "value1")
	storage.Set("key2", "value2")
	storage.Set("key3", "value3")

	seen := make(map[string]string)
	for pair := range storage.Generate() {
		seen[pair[0]] = pair[1]
	}

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	if len(seen) != len(expected) {
		t.Errorf("expected %d elements, got %d", len(expected), len(seen))
	}

	for k, v := range expected {
		if seen[k] != v {
			t.Errorf("expected key %s to have value %s, got %s", k, v, seen[k])
		}
	}
}

func TestMemoryStorageEngineGenerateEmpty(t *testing.T) {
	storage := NewMemoryStorageEngine()
	count := 0
	for range storage.Generate() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 elements, got %d", count)
	}
}
