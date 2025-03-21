package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBRepository struct {
	db *sql.DB
	c  *configs.ServerConfig
}

// NewMetricDBRepository creates a new instance of MetricDBRepository.
func NewMetricDBRepository(c *configs.ServerConfig, db *sql.DB) *MetricDBRepository {
	createMetricsTable(db)
	return &MetricDBRepository{
		db: db,
		c:  c,
	}
}

// SaveMetrics saves a list of metrics in the database.
func (mr *MetricDBRepository) SaveMetrics(ctx context.Context, metrics []*types.Metrics) error {
	// Prepare a query string for bulk insert with ON CONFLICT DO UPDATE.
	query := `INSERT INTO metrics (id, type, delta, value) 
			  VALUES `
	var args []interface{}
	for i, metric := range metrics {
		// Append placeholders and args for each metric
		args = append(args, metric.ID, metric.Type)
		if metric.Delta != nil {
			args = append(args, *metric.Delta)
		} else {
			args = append(args, nil)
		}
		if metric.Value != nil {
			args = append(args, *metric.Value)
		} else {
			args = append(args, nil)
		}

		// Add placeholders for values
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		if i < len(metrics)-1 {
			query += ", "
		}
	}

	// Add ON CONFLICT DO UPDATE clause
	query += ` ON CONFLICT (id, type) 
			   DO UPDATE 
			   SET delta = EXCLUDED.delta, value = EXCLUDED.value`

	// Execute the query
	_, err := mr.db.ExecContext(ctx, query, args...) // Use ExecContext for execute queries with no result rows
	return err
}

// FilterMetricsByTypeAndID filters metrics by their IDs and types, and returns matching metrics.
func (mr *MetricDBRepository) FilterMetricsByTypeAndID(ctx context.Context, metricIDs []types.MetricID) ([]*types.Metrics, error) {
	// Build query with WHERE clause for metric_id and metric_type.
	query := "SELECT id, type, delta, value FROM metrics WHERE "
	var args []interface{}
	for i, metricID := range metricIDs {
		// Add conditions for each metricID (id and type).
		if i > 0 {
			query += " OR "
		}
		query += fmt.Sprintf("(id = $%d AND type = $%d)", i*2+1, i*2+2)
		args = append(args, metricID.ID, metricID.Type)
	}

	// Execute the query and scan results into a slice of Metrics.
	rows, err := mr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %v", err)
	}
	defer rows.Close()

	var metrics []*types.Metrics
	for rows.Next() {
		var metric types.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %v", err)
		}
		metrics = append(metrics, &metric)
	}

	// Handle any row iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %v", err)
	}

	return metrics, nil
}

// ListMetrics lists all metrics stored in the database.
func (mr *MetricDBRepository) ListMetrics(ctx context.Context) ([]*types.Metrics, error) {
	query := "SELECT id, type, delta, value FROM metrics"
	rows, err := mr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %v", err)
	}
	defer rows.Close()

	var metrics []*types.Metrics
	for rows.Next() {
		var metric types.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %v", err)
		}
		metrics = append(metrics, &metric)
	}

	// Handle any row iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %v", err)
	}

	return metrics, nil
}

func createMetricsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		PRIMARY KEY (id, type)
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
