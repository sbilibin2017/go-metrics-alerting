package storage

// Getter управляет операцией чтения из хранилища.
type Getter struct {
	storage *Storage
}

func NewGetter(storage *Storage) *Getter {
	return &Getter{storage: storage}
}

// Get возвращает значение по ключу и флаг, существует ли ключ в хранилище.
func (g *Getter) Get(key string) (string, bool) {
	g.storage.mu.RLock()
	defer g.storage.mu.RUnlock()
	val, exists := g.storage.data[key]
	return val, exists
}
