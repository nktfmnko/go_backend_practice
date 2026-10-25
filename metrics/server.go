package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Общее количество обработанных HTTP-запросов.",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Длительность выполнения HTTP-запросов в секундах.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	requestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "requests_in_flight",
			Help: "Количество HTTP-запросов, обрабатываемых в данный момент.",
		},
	)

	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Общее количество обработанных gRPC-запросов.",
		},
		[]string{"method", "code"},
	)

	grpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Длительность выполнения gRPC-запросов в секундах.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	grpcRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "grpc_requests_in_flight",
			Help: "Количество gRPC-запросов, обрабатываемых в данный момент.",
		},
	)
)

func IncRequests(method, path, status string) {
	requestsTotal.WithLabelValues(method, path, status).Inc()
}

func ObserveDuration(method, path string, duration float64) {
	requestDuration.WithLabelValues(method, path).Observe(duration)
}

func IncInFlight() {
	requestsInFlight.Inc()
}

func DecInFlight() {
	requestsInFlight.Dec()
}

func IncGRPCRequests(method, code string) {
	grpcRequestsTotal.WithLabelValues(method, code).Inc()
}

func ObserveGRPCDuration(method string, duration float64) {
	grpcRequestDuration.WithLabelValues(method).Observe(duration)
}

func IncGRPCInFlight() {
	grpcRequestsInFlight.Inc()
}

func DecGRPCInFlight() {
	grpcRequestsInFlight.Dec()
}

func init() {
	prometheus.MustRegister(
		requestsTotal,
		requestDuration,
		requestsInFlight,
		grpcRequestsTotal,
		grpcRequestDuration,
		grpcRequestsInFlight,
	)
}

func StartMetricsServer(address string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(address, mux)
}
