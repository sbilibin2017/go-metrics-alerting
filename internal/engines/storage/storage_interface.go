package storage

// StorageEngineInterface defines methods for interacting with storage.
type StorageEngineInterface[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Generate() <-chan [2]V
}
