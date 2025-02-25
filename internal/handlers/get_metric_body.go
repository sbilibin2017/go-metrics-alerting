package handlers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/mailru/easyjson"
)

// UpdateMetricService интерфейс для обновления метрики
type GetMetricBodyService interface {
	GetMetric(id string, mtype domain.MetricType) (*domain.Metrics, error)
}

// GetMetricPathHandler создает обработчик с переданным сервисом
func GetMetricBodyHandler(service GetMetricBodyService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMetricRequest

		if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := req.Validate(); err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}

		metric, err := service.GetMetric(req.ID, domain.MetricType(req.MType))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := types.GetMetricBodyResponse{
			GetMetricRequest: types.GetMetricRequest{
				ID:    metric.ID,
				MType: string(metric.MType),
			},
			Delta: metric.Delta,
			Value: metric.Value,
		}

		w.WriteHeader(http.StatusOK)
		data, _ := easyjson.Marshal(resp)
		w.Write(data)
	}
}
