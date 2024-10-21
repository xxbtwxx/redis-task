package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

type Metrics struct {
	processingStats *prometheus.HistogramVec
}

var metrics *Metrics

func ObserveProcessingTimes(consumerID, state string, processingTime float64) {
	histogram, err := metrics.processingStats.GetMetricWithLabelValues(consumerID, state)
	if err != nil {
		log.Error().Err(err).Msg("Failed finding counter with given labels")
		return
	}
	histogram.Observe(processingTime)
}

func setup() {
	metrics = &Metrics{
		processingStats: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "redis_task",
				Subsystem: "processor",
				Name:      "message_process_time",
				Help:      "Processing times",
			},
			[]string{"consumer_id", "state"},
		),
	}
}
