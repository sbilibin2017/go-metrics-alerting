package facades

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"strings"

	"github.com/go-resty/resty/v2"
)

// MetricFacade структура для отправки метрик через HTTP
type MetricFacade struct {
	client *resty.Client
	config *configs.AgentConfig
}

// NewMetricFacade конструктор с зависимостями
func NewMetricFacade(client *resty.Client, config *configs.AgentConfig) *MetricFacade {
	address := config.Address
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	config.Address = address
	return &MetricFacade{
		client: client,
		config: config,
	}
}

// UpdateMetric обновляет метрику, отправляя ее на сервер
func (s *MetricFacade) UpdateMetric(metric *types.UpdateMetricBodyRequest) {
	s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metric).
		Post(s.config.Address + "/update/")
}
