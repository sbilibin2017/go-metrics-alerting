package handlers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/go-chi/chi"
)

// UpdateMetricService интерфейс для обновления метрики
type UpdateMetricPathService interface {
	UpdateMetric(metric *domain.Metrics) (*domain.Metrics, error)
}

// UpdateMetricPathHandlerFactory создает обработчик с переданным сервисом
func UpdateMetricPathHandler(service UpdateMetricPathService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		mtype := chi.URLParam(r, "type")
		value := chi.URLParam(r, "value")

		req := &types.UpdateMetricPathRequest{
			ID:    id,
			MType: mtype,
			Value: value,
		}

		if err := req.Validate(); err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}

		_, err := service.UpdateMetric(req.ToDomain())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Updated metric successfully"))
	}
}
