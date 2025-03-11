package repositories

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-metrics-alerting/internal/engines/db"
	"go-metrics-alerting/internal/engines/file"
	"go-metrics-alerting/internal/types"
	"strings"
	"sync"
)

// SaveMetricRepository is an interface for saving a single metric
type SaveMetricRepository interface {
	SaveMetric(ctx context.Context, metric *types.Metrics) bool
}

// SaveMetricsRepository is an interface for saving multiple metrics
type SaveMetricsRepository interface {
	SaveMetrics(ctx context.Context, metrics []*types.Metrics) bool
}

// GetMetricByIdRepository is an interface for retrieving a metric by its ID
type GetMetricByIdRepository interface {
	GetMetricByID(ctx context.Context, id types.MetricID) *types.Metrics
}

// GetMetricsByIdsRepository is an interface for retrieving multiple metrics by their IDs
type GetMetricsByIdsRepository interface {
	GetMetricsByIDs(ctx context.Context, ids []types.MetricID) map[types.MetricID]*types.Metrics
}

// ListMetricsRepository is an interface for listing all metrics
type ListMetricsRepository interface {
	ListMetrics(ctx context.Context) []*types.Metrics
}

type MetricRepository interface {
	SaveMetricRepository
	SaveMetricsRepository
	GetMetricByIdRepository
	GetMetricsByIdsRepository
	ListMetricsRepository
}

// NewMetricRepository создает основной и файловый репозитории в зависимости от конфигурации.
func NewMetricRepository(db *db.DB, file *file.File) (MetricRepository, MetricRepository) {
	dbRepo := NewDBRepository(db)
	fileRepo := NewFileRepository(file)
	var mainRepo MetricRepository

	if dbRepo != nil {
		mainRepo = dbRepo
	} else if fileRepo != nil {
		mainRepo = fileRepo
	} else {
		mainRepo = NewMemoryRepository()
	}

	if mainRepo == nil {
		return nil, nil
	}

	return mainRepo, fileRepo
}

type MetricDBRepository struct {
	SaveMetricDBRepository
	SaveMetricsDBRepository
	GetMetricByIdDBRepository
	GetMetricsByIdsDBRepository
	ListMetricsDBRepository
}

type MetricFileRepository struct {
	SaveMetricFileRepository
	SaveMetricsFileRepository
	GetMetricByIdFileRepository
	GetMetricsByIdsFileRepository
	ListMetricsFileRepository
}

// MetricMemoryRepository - composed repository that includes all individual repositories.
type MetricMemoryRepository struct {
	SaveMetricMemoryRepository
	SaveMetricsMemoryRepository
	GetMetricByIdMemoryRepository
	GetMetricsByIdsMemoryRepository
	ListMetricsMemoryRepository
}

// NewDBRepository создает репозиторий для работы с базой данных.
func NewDBRepository(db *db.DB) *MetricDBRepository {
	if db == nil {
		return nil
	}
	createMetricsTable(db)
	return &MetricDBRepository{
		SaveMetricDBRepository:      SaveMetricDBRepository{db},
		SaveMetricsDBRepository:     SaveMetricsDBRepository{db},
		GetMetricByIdDBRepository:   GetMetricByIdDBRepository{db},
		GetMetricsByIdsDBRepository: GetMetricsByIdsDBRepository{db},
		ListMetricsDBRepository:     ListMetricsDBRepository{db},
	}
}

// NewFileRepository создает репозиторий для работы с файлом.
func NewFileRepository(file *file.File) *MetricFileRepository {
	if file == nil {
		return nil
	}
	return &MetricFileRepository{
		SaveMetricFileRepository:      SaveMetricFileRepository{file},
		SaveMetricsFileRepository:     SaveMetricsFileRepository{&SaveMetricFileRepository{file}},
		GetMetricByIdFileRepository:   GetMetricByIdFileRepository{file},
		GetMetricsByIdsFileRepository: GetMetricsByIdsFileRepository{&GetMetricByIdFileRepository{file}},
		ListMetricsFileRepository:     ListMetricsFileRepository{file},
	}
}

// NewMemoryRepository создает репозиторий для хранения данных в памяти.
func NewMemoryRepository() *MetricMemoryRepository {
	s := &storage{
		mu:   sync.Mutex{},
		data: make(map[types.MetricID]*types.Metrics),
	}
	return &MetricMemoryRepository{
		SaveMetricMemoryRepository:      SaveMetricMemoryRepository{s},
		SaveMetricsMemoryRepository:     SaveMetricsMemoryRepository{s},
		GetMetricByIdMemoryRepository:   GetMetricByIdMemoryRepository{s},
		GetMetricsByIdsMemoryRepository: GetMetricsByIdsMemoryRepository{s},
		ListMetricsMemoryRepository:     ListMetricsMemoryRepository{s},
	}
}

func createMetricsTable(db *db.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS metrics (
		id TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('counter', 'gauge')),
		delta BIGINT,
		value DOUBLE PRECISION,
		PRIMARY KEY (id, type)
	);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

type SaveMetricDBRepository struct {
	*db.DB
}

func (r *SaveMetricDBRepository) SaveMetric(
	ctx context.Context, metric *types.Metrics,
) bool {
	query := `
		INSERT INTO metrics (id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, type) 
		DO UPDATE SET 
			delta = EXCLUDED.delta, 
			value = EXCLUDED.value
	`

	_, err := r.ExecContext(
		ctx,
		query,
		metric.ID,
		metric.Type,
		metric.Delta,
		metric.Value,
	)

	return err == nil
}

type SaveMetricsDBRepository struct {
	*db.DB
}

func (r *SaveMetricDBRepository) SaveMetrics(
	ctx context.Context, metrics []*types.Metrics,
) bool {
	query := `
		INSERT INTO metrics (id, type, delta, value)
		VALUES %s
		ON CONFLICT (id, type)
		DO UPDATE SET 
			delta = EXCLUDED.delta, 
			value = EXCLUDED.value
	`
	var values []interface{}
	var placeholders []string
	for i, metric := range metrics {
		startIdx := i * 4
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", startIdx+1, startIdx+2, startIdx+3, startIdx+4))
		values = append(values, metric.ID, metric.Type, metric.Delta, metric.Value)
	}

	query = strings.Replace(query, "%s", strings.Join(placeholders, ", "), 1)

	_, err := r.ExecContext(ctx, query, values...)

	return err == nil
}

type GetMetricByIdDBRepository struct {
	*db.DB
}

// GetMetricByID retrieves a single metric by its ID from the database.
func (r *GetMetricByIdDBRepository) GetMetricByID(
	ctx context.Context, id types.MetricID,
) *types.Metrics {
	query := "SELECT id, type, delta, value FROM metrics WHERE id = $1 AND type = $2"
	row := r.QueryRowContext(ctx, query, id.ID, id.Type)

	var metric types.Metrics
	if err := row.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}

	return &metric
}

type GetMetricsByIdsDBRepository struct {
	*db.DB
}

// GetMetricsByIDs retrieves multiple metrics by their IDs from the database.
func (r *GetMetricsByIdsDBRepository) GetMetricsByIDs(
	ctx context.Context, ids []types.MetricID,
) map[types.MetricID]*types.Metrics {
	query := "SELECT id, type, delta, value FROM metrics WHERE (id, type) IN ("
	args := []interface{}{}
	placeholders := []string{}
	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, id.ID, id.Type)
	}
	query += strings.Join(placeholders, ", ") + ")"

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	metricsMap := make(map[types.MetricID]*types.Metrics)
	for _, id := range ids {
		metricsMap[id] = nil
	}
	for rows.Next() {
		var metric types.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			continue
		}
		metricID := types.MetricID{ID: metric.ID, Type: metric.Type}
		metricsMap[metricID] = &metric
	}

	err = rows.Err()
	if err != nil {
		return nil
	}

	return metricsMap
}

type ListMetricsDBRepository struct {
	*db.DB
}

// ListMetrics retrieves all metrics from the database.
func (r *ListMetricsDBRepository) ListMetrics(
	ctx context.Context,
) []*types.Metrics {
	query := "SELECT id, type, delta, value FROM metrics"
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var metrics []*types.Metrics

	for rows.Next() {
		var metric types.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			continue
		}
		metrics = append(metrics, &metric)
	}

	err = rows.Err()
	if err != nil {
		return nil
	}

	return metrics
}

type SaveMetricFileRepository struct {
	*file.File
}

func (r *SaveMetricFileRepository) SaveMetric(
	ctx context.Context, metric *types.Metrics,
) bool {
	// Работаем с уже открытым файлом
	file := r.File
	if file == nil || file.File == nil {
		return false
	}

	// Сканируем файл по строкам
	var updatedLines []string
	scanner := bufio.NewScanner(file.File)
	found := false

	// Проходим по всем строкам в файле
	for scanner.Scan() {
		line := scanner.Text()
		var existingMetric types.Metrics
		if err := json.Unmarshal([]byte(line), &existingMetric); err != nil {
			// Если строка некорректна, добавляем её без изменений
			updatedLines = append(updatedLines, line)
			continue
		}

		// Сравниваем по ID метрики
		if existingMetric.MetricID.ID == metric.MetricID.ID {
			// Обновление существующей метрики
			data, err := json.Marshal(metric)
			if err != nil {
				return false
			}
			updatedLines = append(updatedLines, string(data))
			found = true
		} else {
			// Сохранение старой метрики
			updatedLines = append(updatedLines, line)
		}
	}

	// Если метрика не найдена, добавляем новую
	if !found {
		data, err := json.Marshal(metric)
		if err != nil {
			return false
		}
		updatedLines = append(updatedLines, string(data))
	}

	// Перезаписываем файл с новыми данными
	file.Truncate(0) // Очищаем файл
	file.Seek(0, 0)  // Перемещаемся в начало
	writer := bufio.NewWriter(file.File)

	// Записываем все обновленные строки в файл
	for _, line := range updatedLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return false
		}
	}

	// Завершаем запись
	if err := writer.Flush(); err != nil {
		return false
	}

	return true
}

type SaveMetricsFileRepository struct {
	*SaveMetricFileRepository
}

func (r *SaveMetricsFileRepository) SaveMetrics(
	ctx context.Context, metrics []*types.Metrics,
) bool {
	if r.SaveMetricFileRepository == nil {
		return false
	}

	success := true
	for _, metric := range metrics {
		if !r.SaveMetric(ctx, metric) {
			success = false
		}
	}

	return success
}

type GetMetricByIdFileRepository struct {
	*file.File
}

// GetMetricByID retrieves a single metric from the file.
func (r *GetMetricByIdFileRepository) GetMetricByID(
	ctx context.Context, id types.MetricID,
) *types.Metrics {
	if r.File == nil || r.File.File == nil { // Ensure that File and File.File are not nil
		return nil
	}

	// Reset the read pointer to the beginning of the file
	_, err := r.File.Seek(0, 0)
	if err != nil {
		return nil
	}

	// Create a new JSON decoder using the *os.File as the io.Reader
	decoder := json.NewDecoder(r.File.File) // Use r.File.File, which is the *os.File

	for decoder.More() {
		var m types.Metrics
		if err := decoder.Decode(&m); err != nil {
			return nil
		}
		if m.MetricID.ID == id.ID && m.MetricID.Type == id.Type {
			return &m
		}
	}

	return nil
}

type GetMetricsByIdsFileRepository struct {
	*GetMetricByIdFileRepository
}

// GetMetricsByIDs retrieves multiple metrics from the file.
func (r *GetMetricsByIdsFileRepository) GetMetricsByIDs(
	ctx context.Context, ids []types.MetricID,
) map[types.MetricID]*types.Metrics {
	metricsMap := make(map[types.MetricID]*types.Metrics)
	for _, id := range ids {
		if metric := r.GetMetricByID(ctx, id); metric != nil {
			metricsMap[id] = metric
		}
	}
	return metricsMap
}

type ListMetricsFileRepository struct {
	*file.File
}

// ListMetrics lists all metrics from the file.
func (r *ListMetricsFileRepository) ListMetrics(
	ctx context.Context,
) []*types.Metrics {
	if r.File == nil {
		return nil
	}

	_, err := r.File.Seek(0, 0)
	if err != nil {
		return nil
	}

	var metrics []*types.Metrics
	decoder := json.NewDecoder(r.File.File)
	for decoder.More() {
		var m types.Metrics
		if err := decoder.Decode(&m); err != nil {
			continue
		}
		metrics = append(metrics, &m)
	}

	return metrics
}

type storage struct {
	mu   sync.Mutex
	data map[types.MetricID]*types.Metrics
}

type SaveMetricMemoryRepository struct {
	*storage
}

// SaveMetric saves a single metric in memory.
func (r *SaveMetricMemoryRepository) SaveMetric(
	ctx context.Context, metric *types.Metrics,
) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[metric.MetricID] = metric

	return true
}

type SaveMetricsMemoryRepository struct {
	*storage
}

// SaveMetrics saves multiple metrics in memory.
func (r *SaveMetricsMemoryRepository) SaveMetrics(
	ctx context.Context, metrics []*types.Metrics,
) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, metric := range metrics {
		r.data[metric.MetricID] = metric
	}
	return true
}

type GetMetricByIdMemoryRepository struct {
	*storage
}

// GetMetricByID retrieves a single metric from memory.
func (r *GetMetricByIdMemoryRepository) GetMetricByID(
	ctx context.Context, id types.MetricID,
) *types.Metrics {
	r.mu.Lock()
	defer r.mu.Unlock()

	metric, found := r.data[id]
	if !found {
		return nil
	}

	return metric
}

type GetMetricsByIdsMemoryRepository struct {
	*storage
}

// GetMetricsByIDs retrieves multiple metrics from memory.
func (r *GetMetricsByIdsMemoryRepository) GetMetricsByIDs(
	ctx context.Context, ids []types.MetricID,
) map[types.MetricID]*types.Metrics {
	r.mu.Lock()
	defer r.mu.Unlock()

	metricsMap := make(map[types.MetricID]*types.Metrics)
	for _, id := range ids {
		if metric, found := r.data[id]; found {
			metricsMap[id] = metric
		} else {
			metricsMap[id] = nil
		}
	}

	return metricsMap
}

type ListMetricsMemoryRepository struct {
	*storage
}

// ListMetrics lists all metrics from memory.
func (r *ListMetricsMemoryRepository) ListMetrics(
	ctx context.Context,
) []*types.Metrics {
	r.mu.Lock()
	defer r.mu.Unlock()

	var metrics []*types.Metrics
	for _, metric := range r.data {
		metrics = append(metrics, metric)
	}

	return metrics
}
