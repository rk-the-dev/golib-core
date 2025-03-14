package metricshelper

import "fmt"

// MockMetricsHelper is a mock implementation of MetricsHelper
type MockMetricsHelper struct{}

// NewMockMetricsHelper initializes a mock MetricsHelper
func NewMockMetricsHelper() MetricsHelper {
	return &MockMetricsHelper{}
}

// IncrementCounter (mock) simulates increasing a counter
func (m *MockMetricsHelper) IncrementCounter(name string, labels map[string]string) {
	fmt.Println("📊 [Mock] Incremented Counter:", name, labels)
}

// ObserveHistogram (mock) simulates recording a histogram value
func (m *MockMetricsHelper) ObserveHistogram(name string, value float64, labels map[string]string) {
	fmt.Println("📊 [Mock] Observed Histogram:", name, value, labels)
}

// ObserveSummary (mock) simulates recording a summary value
func (m *MockMetricsHelper) ObserveSummary(name string, value float64, labels map[string]string) {
	fmt.Println("📊 [Mock] Observed Summary:", name, value, labels)
}

// SetGauge (mock) simulates setting a gauge value
func (m *MockMetricsHelper) SetGauge(name string, value float64, labels map[string]string) {
	fmt.Println("📊 [Mock] Set Gauge:", name, value, labels)
}

// StartMetricsServer (mock) simulates starting a metrics server
func (m *MockMetricsHelper) StartMetricsServer(port string) {
	fmt.Println("🚀 [Mock] Metrics server started on port:", port)
}

// Close (mock) returns no error
func (m *MockMetricsHelper) Close() error {
	return nil
}
