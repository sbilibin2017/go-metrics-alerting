package facades

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// MetricFacade структура для отправки метрик через HTTP
type MetricFacade struct {
	client *resty.Client
	config *configs.AgentConfig
	logger *zap.Logger
}

// NewMetricFacade конструктор с зависимостями
func NewMetricFacade(client *resty.Client, config *configs.AgentConfig, logger *zap.Logger) *MetricFacade {
	address := config.Address
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	config.Address = address
	return &MetricFacade{
		client: client,
		config: config,
		logger: logger,
	}
}

// UpdateMetrics метод для обновления метрик
func (s *MetricFacade) UpdateMetric(metric *domain.Metrics) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metric).
		Post(s.config.Address + "/update/")
	if err != nil || resp.StatusCode() >= 400 {
		s.logger.Error("Error sending metric", zap.Error(err), zap.String("metricID", metric.ID))
	} else {
		s.logger.Info("Metric sent successfully", zap.String("metricID", metric.ID))
	}

}
