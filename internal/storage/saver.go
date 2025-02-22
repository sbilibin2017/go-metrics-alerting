package storage

// Saver управляет операцией записи в хранилище.
type Saver struct {
	storage *Storage
}

func NewSaver(storage *Storage) *Saver {
	return &Saver{storage: storage}
}

func (s *Saver) Save(key string, value string) bool {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.data[key] = value
	return true
}
