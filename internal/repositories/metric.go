package repositories

import (
	"context"
	"database/sql"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"os"
)

// MetricRepository holds the three repositories.
type MetricRepository struct {
	DbRepo     *MetricDBRepository
	FileRepo   *MetricFileRepository
	MemoryRepo *MetricMemoryRepository
}

// NewMetricRepository creates a new instance of MetricRepository, containing all three repositories.
func NewMetricRepository(c *configs.ServerConfig, file *os.File, db *sql.DB) *MetricRepository {
	// Initialize each repository
	var dbRepo *MetricDBRepository
	var fileRepo *MetricFileRepository
	var err error

	// Initialize DB repository if DatabaseDSN is provided
	if c.DatabaseDSN != "" {
		dbRepo = NewMetricDBRepository(c, db)
	}

	// Initialize File repository if FileStoragePath is provided
	if c.FileStoragePath != "" {
		fileRepo, err = NewMetricFileRepository(c)
		if err != nil {
			return nil
		}
	}

	// Initialize Memory repository by default
	memoryRepo := NewMetricMemoryRepository()

	// Return the MetricRepository that holds them
	return &MetricRepository{
		DbRepo:     dbRepo,
		FileRepo:   fileRepo,
		MemoryRepo: memoryRepo,
	}
}

// MetricRepositoryInterface defines the common methods for all repositories.
type MetricRepo interface {
	SaveMetrics(ctx context.Context, metrics []*types.Metrics) error
	FilterMetricsByTypeAndId(ctx context.Context, metricIDs []types.MetricID) ([]*types.Metrics, error)
	ListMetrics(ctx context.Context) ([]*types.Metrics, error)
}

// GetMainRepository returns the repository with the highest priority (db -> file -> memory), based on ServerConfig.
func (mr *MetricRepository) GetMainRepository(c *configs.ServerConfig) MetricRepo {
	// Check DB repository if DatabaseDSN is set
	if mr.DbRepo != nil && c.DatabaseDSN != "" {
		return mr.DbRepo
	}

	// Check File repository if FileStoragePath is set
	if mr.FileRepo != nil && c.FileStoragePath != "" {
		return mr.FileRepo
	}

	// Fall back to Memory repository if both DB and File are unavailable
	if mr.MemoryRepo != nil {
		return mr.MemoryRepo
	}

	// If no repository is available, return nil or handle the error
	return nil
}
