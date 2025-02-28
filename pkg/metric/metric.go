package metric

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}
type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func CreateMetrics(address string, name string) (Metrics, error) {

	metr := &PrometheusMetrics{}

	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
		Help: "Total number of hits",
	})

	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
			Help: "Total number of hits",
		},
		[]string{"status", "method", "path"},
	)
	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name + "_response_time",
			Help:    "Response time histogram",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "method", "path"},
	)

	// Đăng ký metrics với Prometheus
	for _, metric := range []prometheus.Collector{metr.HitsTotal, metr.Hits, metr.Times} {
		if err := prometheus.Register(metric); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, err
			}
		}
	}

	go func() {
		router := gin.Default()
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))

		// Lấy cổng từ address
		port := "unknown"
		if idx := strings.LastIndex(address, ":"); idx != -1 {
			port = address[idx+1:]
		}

		log.Printf("Metrics server is running on port: %s", port)
		if err := router.Run(address); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	return metr, nil
}

func (metr *PrometheusMetrics) IncHits(status int, method, path string) {
	metr.HitsTotal.Inc()
	metr.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metr *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metr.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}
