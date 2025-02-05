package facades

import (
	"fmt"
	"go-metrics-alerting/internal/apiclient"
	"go-metrics-alerting/pkg/logger"
	"net/http"
)

// PostRequester интерфейс для выполнения POST-запросов
type PostRequester interface {
	Post(url string, headers map[string]string) (*apiclient.APIResponse, error)
}

// Ошибки, связанные с обновлением метрик
var (
	ErrInvalidStatusCode = fmt.Errorf("invalid status code while updating metric")
	ErrRequestFailed     = fmt.Errorf("failed to send the request")
)

// MetricFacade фасад для работы с метриками
type metricFacade struct {
	client  PostRequester
	address string
}

// NewMetricFacade создает новый экземпляр MetricFacade
func NewMetricFacade(client PostRequester, address string) *metricFacade {
	return &metricFacade{
		client:  client,
		address: address,
	}
}

// Update обновляет значение метрики в зависимости от типа метрики
func (mf *metricFacade) UpdateMetric(metricType string, metricName string, metricValue string) error {
	// Формирование URL
	url := fmt.Sprintf("%s/update/%s/%s/%s", mf.address, metricType, metricName, metricValue)

	// Логируем отправку запроса
	logger.Logger.Debugf("Sending POST request to URL: %s", url)

	// Отправляем POST-запрос с помощью PostRequester
	resp, err := mf.client.Post(url, map[string]string{
		"Content-Type": "text/plain",
	})
	if err != nil {
		// Логируем ошибку при отправке запроса
		logger.Logger.Errorf("Error sending POST request to %s: %v", url, err)
		return ErrRequestFailed
	}

	// Логируем статус ответа
	logger.Logger.Debugf("Received response: %d %s", resp.StatusCode, url)

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		// Логируем ошибку, если статус код не 200
		logger.Logger.Errorf("Invalid status code %d received from %s: %s", resp.StatusCode, url, resp.Body)
		return ErrInvalidStatusCode
	}

	// Логируем успешное обновление
	logger.Logger.Infof("Successfully updated metric %s %s = %s", metricType, metricName, metricValue)

	return nil
}
