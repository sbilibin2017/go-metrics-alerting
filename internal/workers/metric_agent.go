package workers

import (
	"time"

	"github.com/sbilibin2017/go-metrics-alerting/internal/configs"
	"github.com/sbilibin2017/go-metrics-alerting/internal/services"
)

type MetricsWorkerInterface interface {
	Start()
}

type MetricsWorker struct {
	agentService services.MetricAgentServiceInterface
	config       configs.AgentConfigInterface
}

func NewMetricsWorker(
	agentService services.MetricAgentServiceInterface,
	config configs.AgentConfigInterface,
) MetricsWorkerInterface {
	return &MetricsWorker{
		agentService: agentService,
		config:       config,
	}
}

func (w *MetricsWorker) Start() {
	pollTicker := time.NewTicker(time.Duration(w.config.GetPollInterval()))
	reportTicker := time.NewTicker(time.Duration(w.config.GetReportInterval()))

	go func() {
		for {
			select {
			case <-pollTicker.C:
				w.agentService.CollectMetrics()
			case <-reportTicker.C:
			}
		}
	}()

	<-make(chan struct{})
}
