package metricshelper

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHelper defines the interface for exporting metrics
type MetricsHelper interface {
	IncrementCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	ObserveSummary(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
	StartMetricsServer(port string)
	Close() error
}

// metricsHelper implements MetricsHelper
type metricsHelper struct {
	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	histograms map[string]*prometheus.HistogramVec
	summaries  map[string]*prometheus.SummaryVec
}

var (
	instance *metricsHelper
	once     sync.Once
)

// NewMetricsHelper initializes and returns a MetricsHelper instance
func NewMetricsHelper() MetricsHelper {
	once.Do(func() {
		instance = &metricsHelper{
			counters:   make(map[string]*prometheus.CounterVec),
			gauges:     make(map[string]*prometheus.GaugeVec),
			histograms: make(map[string]*prometheus.HistogramVec),
			summaries:  make(map[string]*prometheus.SummaryVec),
		}
	})
	return instance
}

// IncrementCounter increases the counter metric by 1
func (m *metricsHelper) IncrementCounter(name string, labels map[string]string) {
	if _, exists := m.counters[name]; !exists {
		m.counters[name] = prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: "Counter for " + name}, getLabelKeys(labels))
		prometheus.MustRegister(m.counters[name])
	}
	m.counters[name].With(labels).Inc()
}

// ObserveHistogram records a value in a histogram
func (m *metricsHelper) ObserveHistogram(name string, value float64, labels map[string]string) {
	if _, exists := m.histograms[name]; !exists {
		m.histograms[name] = prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: name, Help: "Histogram for " + name, Buckets: prometheus.DefBuckets}, getLabelKeys(labels))
		prometheus.MustRegister(m.histograms[name])
	}
	m.histograms[name].With(labels).Observe(value)
}

// ObserveSummary records a value in a summary
func (m *metricsHelper) ObserveSummary(name string, value float64, labels map[string]string) {
	if _, exists := m.summaries[name]; !exists {
		m.summaries[name] = prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: name, Help: "Summary for " + name}, getLabelKeys(labels))
		prometheus.MustRegister(m.summaries[name])
	}
	m.summaries[name].With(labels).Observe(value)
}

// SetGauge sets a gauge metric to a specific value
func (m *metricsHelper) SetGauge(name string, value float64, labels map[string]string) {
	if _, exists := m.gauges[name]; !exists {
		m.gauges[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: "Gauge for " + name}, getLabelKeys(labels))
		prometheus.MustRegister(m.gauges[name])
	}
	m.gauges[name].With(labels).Set(value)
}

// StartMetricsServer starts an HTTP server to expose Prometheus metrics
func (m *metricsHelper) StartMetricsServer(port string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("🚀 Metrics server running on port:", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println("❌ Error starting metrics server:", err)
		}
	}()
}

// Close is a placeholder to maintain interface consistency
func (m *metricsHelper) Close() error {
	return nil
}

// getLabelKeys extracts keys from a map for Prometheus labels
func getLabelKeys(labels map[string]string) []string {
	keys := []string{}
	for key := range labels {
		keys = append(keys, key)
	}
	return keys
}
