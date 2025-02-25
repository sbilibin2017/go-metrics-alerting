package handlers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/go-chi/chi"
)

// UpdateMetricService интерфейс для обновления метрики
type GetMetricPathService interface {
	GetMetric(id string, mtype domain.MetricType) (*domain.Metrics, error)
}

// GetMetricPathHandler создает обработчик с переданным сервисом
func GetMetricPathHandler(service GetMetricPathService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		mtype := chi.URLParam(r, "type")

		req := &types.GetMetricRequest{
			ID:    id,
			MType: mtype,
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

		w.WriteHeader(http.StatusOK)
		response := types.GetMetricPathResponse{}
		response = response.FromDomain(metric)
		w.Write([]byte(response.Value))
	}
}
