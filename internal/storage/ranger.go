package storage

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger struct {
	storage *Storage
}

func NewRanger(storage *Storage) *Ranger {
	return &Ranger{storage: storage}
}

func (r *Ranger) Range(callback func(key string, value string) bool) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()
	for key, value := range r.storage.data {
		if !callback(key, value) {
			break
		}
	}
}
