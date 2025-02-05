package repositories

// Storage интерфейс для работы с хранилищем данных
type StorageSetter interface {
	Set(key string, value string) error
}

// Storage интерфейс для работы с хранилищем данных
type StorageGetter interface {
	Get(key string) (string, error)
}

// Storage интерфейс для работы с хранилищем данных
type StorageGenerator interface {
	Generate() <-chan [2]string
}

// Storage интерфейс для работы с хранилищем данных
type Storage interface {
	StorageSetter
	StorageGetter
	StorageGenerator
}

// KeyProcessor интерфейс для кодирования и декодирования ключей
type KeyEncoder interface {
	Encode(metricType string, metricName string) string
}

// KeyProcessor интерфейс для кодирования и декодирования ключей
type KeyDecoder interface {
	Decode(key string) (string, string, error)
}

type KeyProcessor interface {
	KeyEncoder
	KeyDecoder
}

// MetricRepository для работы с хранилищем метрик
type MetricRepository struct {
	storage      Storage
	keyProcessor KeyProcessor
}

// NewMetricRepository создает новый экземпляр MetricRepository
func NewMetricRepository(storage Storage, keyProcessor KeyProcessor) *MetricRepository {
	return &MetricRepository{
		storage:      storage,
		keyProcessor: keyProcessor,
	}
}

// Save сохраняет метрику в хранилище
func (r *MetricRepository) Save(metricType string, metricName string, value string) error {
	// Генерация ключа
	keyEncoded := r.keyProcessor.Encode(metricType, metricName)

	// Сохранение в хранилище
	return r.storage.Set(keyEncoded, value)
}

// Get извлекает метрику из хранилища по ключу
func (r *MetricRepository) Get(metricType string, metricName string) (string, error) {
	// Генерация ключа
	key := r.keyProcessor.Encode(metricType, metricName)

	// Извлечение из хранилища
	return r.storage.Get(key)
}

// GetAll извлекает все метрики из хранилища с использованием метода Generate
func (r *MetricRepository) GetAll() [][3]string {
	// Создаем срез для хранения всех элементов
	var allMetrics [][3]string

	// Читаем из канала
	for kv := range r.storage.Generate() {
		// Декодируем ключ для получения типа метрики и имени метрики
		metricType, metricName, err := r.keyProcessor.Decode(kv[0])
		if err != nil {
			continue
		}

		// Добавляем элемент в срез
		allMetrics = append(allMetrics, [3]string{metricType, metricName, kv[1]})
	}

	return allMetrics
}
