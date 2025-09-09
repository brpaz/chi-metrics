package metrics

import "github.com/prometheus/client_golang/prometheus"

// Counter creates a counter metric using the default registry
func Counter(name string, help string) CounterMetric {
	return CounterWithRegistry(nil, name, help)
}

// CounterWith creates a counter metric with typed labels using the default registry
func CounterWith[T any](name string, help string) CounterMetricLabeled[T] {
	return CounterWithRegistryWith[T](nil, name, help)
}

// CounterWithRegistry creates a counter metric using the specified registry.
// If registry is nil, the default registry is used.
func CounterWithRegistry(registry *prometheus.Registry, name string, help string) CounterMetric {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: mustValidMetricName(name),
		Help: help,
	}, []string{})
	
	if registry != nil {
		registry.MustRegister(vec)
	} else {
		prometheus.MustRegister(vec)
	}
	return CounterMetric{Vec: vec}
}

// CounterWithRegistryWith creates a counter metric with typed labels using the specified registry.
// If registry is nil, the default registry is used.
func CounterWithRegistryWith[T any](registry *prometheus.Registry, name string, help string) CounterMetricLabeled[T] {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: mustValidMetricName(name),
		Help: help,
	}, getLabelKeys[T]())
	
	if registry != nil {
		registry.MustRegister(vec)
	} else {
		prometheus.MustRegister(vec)
	}
	return CounterMetricLabeled[T]{Vec: vec}
}

type CounterMetric struct {
	Vec *prometheus.CounterVec
}

func (c *CounterMetric) Inc() {
	c.Vec.With(prometheus.Labels{}).Inc()
}

func (c *CounterMetric) Add(value float64) {
	c.Vec.With(prometheus.Labels{}).Add(value)
}

type CounterMetricLabeled[T any] struct {
	Vec *prometheus.CounterVec
}

func (c *CounterMetricLabeled[T]) Inc(labels T) {
	c.Vec.With(getLabelValues(labels)).Inc()
}

func (c *CounterMetricLabeled[T]) Add(value float64, labels T) {
	c.Vec.With(getLabelValues(labels)).Add(value)
}
