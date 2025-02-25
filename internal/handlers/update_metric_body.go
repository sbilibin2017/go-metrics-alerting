package handlers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/mailru/easyjson"
)

// UpdateMetricBodyService интерфейс для обновления метрики через тело запроса
type UpdateMetricBodyService interface {
	UpdateMetric(metric *domain.Metrics) (*domain.Metrics, error)
}

// Обработчик для обновления метрики через тело запроса с использованием замыкания
func UpdateMetricBodyHandler(service UpdateMetricBodyService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMetricBodyRequest

		if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := req.Validate(); err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}

		metric := req.ToDomain()

		updatedMetric, err := service.UpdateMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := types.UpdateMetricBodyResponse{
			UpdateMetricBodyRequest: types.UpdateMetricBodyRequest{
				ID:    updatedMetric.ID,
				MType: string(updatedMetric.MType),
				Delta: updatedMetric.Delta,
				Value: updatedMetric.Value,
			},
		}

		w.WriteHeader(http.StatusOK)
		data, _ := easyjson.Marshal(resp)
		w.Write(data)
	}
}
