package facades

import (
	"encoding/json"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
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

// UpdateMetric обновляет метрику, отправляя ее на сервер
func (s *MetricFacade) UpdateMetric(metric *types.UpdateMetricBodyRequest) {
	metricBody, _ := json.Marshal(metric)
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metricBody).
		Post(s.config.Address + "/update/")

	if err != nil {
		s.logger.Error("Error sending metric", zap.String("metricID", metric.ID))
		return
	}
	// Логируем успешную отправку
	s.logger.Info("Metric sent successfully", zap.String("metricID", metric.ID), zap.Int("statusCode", resp.StatusCode()))

}
