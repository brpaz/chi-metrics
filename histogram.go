package metrics

import "github.com/prometheus/client_golang/prometheus"

// Histogram creates a histogram metric using the default registry
func Histogram(name, help string, buckets []float64) HistogramMetric {
	return HistogramWithRegistry(nil, name, help, buckets)
}

// HistogramWith creates a histogram metric with typed labels using the default registry
func HistogramWith[T any](name, help string, buckets []float64) HistogramMetricLabeled[T] {
	return HistogramWithRegistryWith[T](nil, name, help, buckets)
}

// HistogramWithRegistry creates a histogram metric using the specified registry.
// If registry is nil, the default registry is used.
func HistogramWithRegistry(registry *prometheus.Registry, name, help string, buckets []float64) HistogramMetric {
	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    mustValidMetricName(name),
		Help:    help,
		Buckets: buckets,
	}, []string{})
	
	if registry != nil {
		registry.MustRegister(vec)
	} else {
		prometheus.MustRegister(vec)
	}
	return HistogramMetric{Vec: vec}
}

// HistogramWithRegistryWith creates a histogram metric with typed labels using the specified registry.
// If registry is nil, the default registry is used.
func HistogramWithRegistryWith[T any](registry *prometheus.Registry, name, help string, buckets []float64) HistogramMetricLabeled[T] {
	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    mustValidMetricName(name),
		Help:    help,
		Buckets: buckets,
	}, getLabelKeys[T]())
	
	if registry != nil {
		registry.MustRegister(vec)
	} else {
		prometheus.MustRegister(vec)
	}
	return HistogramMetricLabeled[T]{Vec: vec}
}

type HistogramMetric struct {
	Vec *prometheus.HistogramVec
}

func (h *HistogramMetric) Observe(value float64) {
	h.Vec.With(prometheus.Labels{}).Observe(value)
}

// HistogramMetric represents a histogram metric with typed labels
type HistogramMetricLabeled[T any] struct {
	Vec *prometheus.HistogramVec
}

func (h *HistogramMetricLabeled[T]) Observe(value float64, labels T) {
	h.Vec.With(getLabelValues(labels)).Observe(value)
}
