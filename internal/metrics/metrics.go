package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TaskManagerMetrics *Metrics

type Metrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestsDuration *prometheus.HistogramVec
	ErrorsTotal      *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		// Счетчик запросов по методу и пути
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "task_manager_requests_total",
				Help: "Total number for requests",
			},
			[]string{"method", "path"},
		),
		// Гистограмма времени выполнения запросов
		RequestsDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "task_manager_request_duration_seconds",
				Help:    "Duration of requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		// Счетчик ошибок
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "task_manager_errors_total",
				Help: "Total number of errors",
			},
			[]string{"method", "path", "status"},
		),
	}
}

func InitPrometheus() {
	TaskManagerMetrics = NewMetrics()
}
