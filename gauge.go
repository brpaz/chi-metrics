package metrics

import "github.com/prometheus/client_golang/prometheus"

// Gauge creates a gauge metric using the default registry
func Gauge(name, help string) GaugeMetric {
	return GaugeWithRegistry(nil, name, help)
}

// GaugeWith creates a gauge metric with typed labels using the default registry
func GaugeWith[T any](name, help string) GaugeMetricLabeled[T] {
	return GaugeWithRegistryWith[T](nil, name, help)
}

// GaugeWithRegistry creates a gauge metric using the specified registry.
// If registry is nil, the default registry is used.
func GaugeWithRegistry(registry *prometheus.Registry, name, help string) GaugeMetric {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: mustValidMetricName(name),
		Help: help,
	}, []string{})
	
	if registry != nil {
		registry.MustRegister(vec)
	} else {
		prometheus.MustRegister(vec)
	}
	return GaugeMetric{Vec: vec}
}

// GaugeWithRegistryWith creates a gauge metric with typed labels using the specified registry.
// If registry is nil, the default registry is used.
func GaugeWithRegistryWith[T any](registry *prometheus.Registry, name, help string) GaugeMetricLabeled[T] {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: mustValidMetricName(name),
		Help: help,
	}, getLabelKeys[T]())
	
	if registry != nil {
		registry.MustRegister(vec)
	} else {
		prometheus.MustRegister(vec)
	}
	return GaugeMetricLabeled[T]{Vec: vec}
}

type GaugeMetric struct {
	Vec *prometheus.GaugeVec
}

func (g *GaugeMetric) Set(value float64) {
	g.Vec.With(prometheus.Labels{}).Set(value)
}

func (g *GaugeMetric) Add(value float64) {
	g.Vec.With(prometheus.Labels{}).Add(value)
}

func (g *GaugeMetric) Inc() {
	g.Vec.With(prometheus.Labels{}).Add(1.0)
}

func (g *GaugeMetric) Dec() {
	g.Vec.With(prometheus.Labels{}).Add(-1.0)
}

type GaugeMetricLabeled[T any] struct {
	Vec *prometheus.GaugeVec
}

func (g *GaugeMetricLabeled[T]) Set(value float64, labels T) {
	g.Vec.With(getLabelValues(labels)).Set(value)
}

func (g *GaugeMetricLabeled[T]) Add(value float64, labels T) {
	g.Vec.With(getLabelValues(labels)).Add(value)
}

func (g *GaugeMetricLabeled[T]) Inc(labels T) {
	g.Vec.With(getLabelValues(labels)).Add(1.0)
}

func (g *GaugeMetricLabeled[T]) Dec(labels T) {
	g.Vec.With(getLabelValues(labels)).Add(-1.0)
}
