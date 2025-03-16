package metric

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	// HTTP metrics
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)

	// Business metrics
	IncActiveUsers(userType string)
	DecActiveUsers(userType string)
	IncMessages(messageType string)

	// System metrics
	SetCPUUsage(label string, usage float64)
	SetMemoryUsage(memType string, usage float64)
}

// PrometheusMetrics implementation
type PrometheusMetrics struct {
	// HTTP metrics
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec

	// Business metrics
	ActiveUsers  *prometheus.GaugeVec
	MessagesSent *prometheus.CounterVec

	// System metrics
	CPUUsage    *prometheus.GaugeVec
	MemoryUsage *prometheus.GaugeVec
}

func NewMetrics(namespace string) (*PrometheusMetrics, error) {
	m := &PrometheusMetrics{
		HitsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "hits_total",
			Help:      "Total number of hits",
		}),
		Hits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "hits",
				Help:      "Total number of hits by status, method, path",
			}, []string{"status", "method", "path"},
		),

		Times: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "response_time",
				Help:      "Response time histogram",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"status", "method", "path"},
		),
		ActiveUsers: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "active_users",
				Help:      "Number of active users",
			},
			[]string{"type"},
		),

		MessagesSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "messages_sent_total",
				Help:      "Total number of messages sent",
			},
			[]string{"type"},
		),

		// System metrics
		CPUUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cpu_usage",
				Help:      "CPU usage percentage",
			},
			[]string{"core"},
		),

		MemoryUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_usage",
				Help:      "Memory usage in bytes",
			},
			[]string{"type"},
		),
	}
	// Register all metrics
	collectors := []prometheus.Collector{
		m.HitsTotal,
		m.Hits,
		m.Times,
		m.ActiveUsers,
		m.MessagesSent,
		m.CPUUsage,
		m.MemoryUsage,
	}
	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, err
			}
		}
	}

	return m, nil
}

// HTTP metrics methods
func (m *PrometheusMetrics) IncHits(status int, method, path string) {
	statusStr := strconv.Itoa(status)
	m.HitsTotal.Inc()
	m.Hits.WithLabelValues(statusStr, method, path).Inc()
}
func (m *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	statusStr := strconv.Itoa(status)
	m.Times.WithLabelValues(statusStr, method, path).Observe(observeTime)
}

// Business metrics methods
func (m *PrometheusMetrics) IncActiveUsers(userType string) {
	m.ActiveUsers.WithLabelValues(userType).Inc()
}

func (m *PrometheusMetrics) DecActiveUsers(userType string) {
	m.ActiveUsers.WithLabelValues(userType).Dec()
}

func (m *PrometheusMetrics) IncMessages(messageType string) {
	m.MessagesSent.WithLabelValues(messageType).Inc()
}

// System metrics methods
func (m *PrometheusMetrics) SetCPUUsage(core string, usage float64) {
	m.CPUUsage.WithLabelValues(core).Set(usage)
}

func (m *PrometheusMetrics) SetMemoryUsage(memType string, usage float64) {
	m.MemoryUsage.WithLabelValues(memType).Set(usage)
}
