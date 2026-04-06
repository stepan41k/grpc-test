package metrics

import (
	"net/http"

	srvmetrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Reg = prometheus.NewRegistry()

	GrpcMetrics = srvmetrics.NewServerMetrics(
		srvmetrics.WithServerHandlingTimeHistogram(
			srvmetrics.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	ExternalAPIRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grinex_api_requests_total",
			Help: "Total number of requests to Grinex API",
		},
		[]string{"status"},
	)

	LastUSDTPrice = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "usdt_current_price",
			Help: "Current USDT price from last successful fetch",
		},
	)
)

func init() {
	Reg.MustRegister(GrpcMetrics)
	Reg.MustRegister(ExternalAPIRequests)
	Reg.MustRegister(LastUSDTPrice)

	Reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	Reg.MustRegister(collectors.NewGoCollector())
}

func StartMetricsServer(addr string) error {
	http.Handle("/metrics", promhttp.HandlerFor(Reg, promhttp.HandlerOpts{}))
	return http.ListenAndServe(addr, nil)
}
