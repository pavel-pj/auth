package metric

import (
	"context"

	"github.com/pavel-pj/auth/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	requestCounter prometheus.Counter
}

var metrics *Metrics

func Init(_ context.Context, cfg config.PrometeusConfig) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace(),
				Subsystem: "grpc",
				Name:      cfg.AppName() + "_request_total",
				Help:      "Количество запросов в секунду",
			},
		),
	}

	return nil

}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
