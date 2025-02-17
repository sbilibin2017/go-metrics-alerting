package facades

import (
	"errors"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var (
	ErrNetwork = errors.New("network error")
	ErrStatus  = errors.New("unexpected status code")
)

// MetricFacade представляет фасад для работы с API метрик.
type MetricFacade struct {
	client *resty.Client
	config *configs.AgentConfig
}

// NewMetricFacade создает новый экземпляр MetricFacade.
func NewMetricFacade(client *resty.Client, config *configs.AgentConfig) *MetricFacade {
	return &MetricFacade{
		client: client,
		config: config,
	}
}

// UpdateMetric отправляет метрику на сервер.
func (c *MetricFacade) UpdateMetric(metric types.MetricsRequest) error {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metric).
		Post(fmt.Sprintf("%s/update/", c.config.Address))

	if err != nil {
		return ErrNetwork
	}

	if resp.StatusCode() != http.StatusOK {
		return ErrStatus
	}

	return nil
}
