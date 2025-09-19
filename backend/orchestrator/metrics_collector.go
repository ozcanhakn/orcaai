package orchestrator

import "time"

// ProviderMetrics represents basic runtime stats per provider used by LB/Brain
type ProviderMetrics struct {
    SuccessRate       float64
    ErrorRate         float64
    TimeoutRate       float64
    ActiveConnections int
    MaxConnections    int
    AverageLatency    time.Duration
}

// MetricsCollector is a lightweight wrapper supplying provider metrics
type MetricsCollector struct{}

func NewMetricsCollector() *MetricsCollector { return &MetricsCollector{} }

func (m *MetricsCollector) GetProviderMetrics(provider string) *ProviderMetrics {
    // Return sane defaults; in real impl, aggregate from Prometheus/DB
    return &ProviderMetrics{
        SuccessRate:       0.98,
        ErrorRate:         0.01,
        TimeoutRate:       0.01,
        ActiveConnections: 1,
        MaxConnections:    100,
        AverageLatency:    800 * time.Millisecond,
    }
}

func (m *MetricsCollector) RecordSuccess(provider string) {}
func (m *MetricsCollector) RecordError(provider string, err error) {}


