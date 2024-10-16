package metrics

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Metrics struct {
	processingStats *prometheus.HistogramVec
}

var metrics *Metrics

func ObserveProcessingTimes(consumerID string, processingTime float64) {
	histogram, err := metrics.processingStats.GetMetricWithLabelValues(consumerID)
	if err != nil {
		log.Error().Err(err).Msg("Failed finding counter with given labels")
		return
	}
	histogram.Observe(processingTime)
}

func Expose() {
	metrics = &Metrics{
		processingStats: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "redis_task",
				Subsystem: "processor",
				Name:      "message_process_time",
				Help:      "Processing times",
			},
			[]string{"consumer_id"},
		),
	}

	go listenAndServe()
}

func listenAndServe() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to start metrics handler")
		os.Exit(1)
	}
}
