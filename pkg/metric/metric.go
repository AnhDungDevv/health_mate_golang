package metric

import (
	"log"
	"strconv"

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

	var metr PrometheusMetrics

	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits-total",
	})

	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, err
	}
	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{
			"status",
			"method",
			"path"},
	)
	if err := prometheus.Register(metr.Hits); err != nil {
		return nil, err
	}
	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.Times); err != nil {
		return nil, err
	}
	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		return nil, err
	}
	go func() {
		router := gin.Default()

		router.GET("/metric", gin.WrapH(promhttp.Handler()))

		log.Printf("Metrics server is running on port: %s", address)
		if err := router.Run(address); err != nil {
			log.Fatal(err)
		}
	}()

	return &metr, nil
}

func (metr *PrometheusMetrics) IncHits(status int, method, path string){
	metr.HitsTotal.Inc()
	metr.Hits.WithLabelValues(strconv.Itoa(status), method,path).Inc()
}

func (metr *PrometheusMetrics) ObserveResponseTime(status int, method, path string , observeTime float64 ){
	metr.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime);
}