package facades

import (
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/types"
)

// MetricFacadeInterface определяет методы для отправки метрик.
type MetricFacadeInterface interface {
	// SendMetric отправляет метрику на сервер и возвращает тело ответа и статус.
	SendMetric(req *types.UpdateMetricRequest) ([]byte, int, error)
}

// MetricFacade предоставляет удобный интерфейс для отправки метрик через ApiClientInterface.
type MetricFacade struct {
	config    configs.AgentConfigInterface
	apiClient engines.ApiClientInterface
}

// NewMetricFacade создает новый экземпляр MetricFacade.
func NewMetricFacade(apiClient engines.ApiClientInterface, config configs.AgentConfigInterface) *MetricFacade {
	return &MetricFacade{
		config:    config,
		apiClient: apiClient,
	}
}

// SendMetric отправляет метрику на сервер.
func (m *MetricFacade) SendMetric(req *types.UpdateMetricRequest) ([]byte, int, error) {
	// Формируем путь с учетом параметров из запроса
	path := fmt.Sprintf("/update/%v/%v/%v", req.Type, req.Name, fmt.Sprintf("%v", req.Value))

	// Используем глобальный resty-клиент из пакета engines
	resp, err := m.apiClient.R().SetHeader("Content-Type", "text/plain").Post(m.config.GetServerURL() + path)

	// Обработка ошибок
	if err != nil {
		return nil, 0, err
	}

	// Возвращаем тело ответа, статус код и nil для ошибки, если все прошло успешно
	return resp.Body(), resp.StatusCode(), nil
}
