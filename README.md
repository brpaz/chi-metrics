# metrics

> Go package for metrics collection in [OpenMetrics](https://github.com/prometheus/OpenMetrics/blob/main/specification/OpenMetrics.md) format.

[![Go Reference](https://pkg.go.dev/badge/github.com/go-chi/metrics.svg)](https://pkg.go.dev/github.com/go-chi/metrics)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-chi/metrics)](https://goreportcard.com/report/github.com/go-chi/httplog)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- **🚀 High Performance**: Built on top of [Prometheus](https://github.com/prometheus/client_golang) Go client with minimal overhead.
- **🌐 HTTP Middleware**: Real-time monitoring of incoming requests.
- **🔄 HTTP Transport**: Client instrumentation for outgoing requests.
- **🎯 Compatibility**: Compatible with [OpenMetrics 1.0](https://github.com/prometheus/OpenMetrics/blob/main/specification/OpenMetrics.md) collectors, e.g. Prometheus.
- **🔒 Type Safety**: Compile-time type-safe metric labels with struct tags validation.
- **🏷️ Data Cardinality**: The API helps you keep the metric label cardinality low.
- **📊 Complete Metrics**: Counter, Gauge, and Histogram metrics with customizable buckets.

## Usage

`go get github.com/go-chi/metrics@latest`

### Basic Usage (Default Registry)

```go
package main

import (
	"github.com/go-chi/metrics"
)

func main() {
	r := chi.NewRouter()

	// Collect metrics for incoming HTTP requests automatically.
	r.Use(metrics.Collector(metrics.CollectorOpts{
		Host:  false,
		Proto: true,
		Skip: func(r *http.Request) bool {
			return r.Method != "OPTIONS"
		},
	}))

	r.Handle("/metrics", metrics.Handler())
	r.Post("/do-work", doWork)

	// Collect metrics for outgoing HTTP requests automatically.
	transport := metrics.Transport(metrics.TransportOpts{
		Host: true,
	})
	http.DefaultClient.Transport = transport(http.DefaultTransport)

	go simulateTraffic()

	log.Println("Server starting on :8022")
	if err := http.ListenAndServe(":8022", r); err != nil {
		log.Fatal(err)
	}
}

// Strongly typed metric labels help maintain low data cardinality
// by enforcing consistent label names across the codebase.
type jobLabels struct {
	Name   string `label:"name"`
	Status string `label:"status"`
}

var jobCounter = metrics.CounterWith[jobLabels]("jobs_processed_total", "Number of jobs processed")

func doWork(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second) // simulate work

	if rand.Intn(100) > 90 { // simulate error
		jobCounter.Inc(jobLabels{Name: "job", Status: "error"})
		w.Write([]byte("Job failed.\n"))
		return
	}

	jobCounter.Inc(jobLabels{Name: "job", Status: "success"})
	w.Write([]byte("Job finished successfully.\n"))
}

func simulateTraffic() {
	for {
		_, _ = client.Get("http://example.com")
		time.Sleep(500 * time.Millisecond)
	}
}
```

### Custom Registry Usage

You can use custom Prometheus registries to isolate metrics or collect them separately:

```go
package main

import (
	"github.com/go-chi/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	// Create a custom registry
	customRegistry := prometheus.NewRegistry()

	r := chi.NewRouter()

	// Use collector with custom registry
	r.Use(metrics.Collector(metrics.CollectorOpts{
		Registry: customRegistry,
		Host:     true,
		Proto:    true,
	}))

	// Serve metrics from default registry (custom registries use Collector only)
	r.Handle("/metrics", metrics.Handler())

	// Create metrics in the custom registry
	counter := metrics.CounterWithRegistry(customRegistry, "app_requests_total", "Total app requests")
	gauge := metrics.GaugeWithRegistry(customRegistry, "app_connections", "Active connections")

	// Typed labeled metrics with custom registry
	type requestLabels struct {
		Method string `label:"method"`
		Status string `label:"status"`
	}
	labeledCounter := metrics.CounterWithRegistryWith[requestLabels](
		customRegistry,
		"app_requests_by_method_total", 
		"Requests by method and status",
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		counter.Inc()
		gauge.Set(42)
		labeledCounter.Inc(requestLabels{Method: "GET", Status: "200"})
		w.Write([]byte("Hello!\n"))
	})

	// Use custom registry for outgoing requests too
	transport := metrics.Transport(metrics.TransportOpts{
		Registry: customRegistry,
		Host:     true,
	})
	http.DefaultClient.Transport = transport(http.DefaultTransport)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
```

### Multiple Registries

You can use different registries for different purposes:

```go
func multipleRegistriesExample() {
	// Business metrics registry
	businessRegistry := prometheus.NewRegistry()
	// Infrastructure metrics registry  
	infraRegistry := prometheus.NewRegistry()

	// Business metrics - use with Collector middleware
	businessCounter := metrics.CounterWithRegistry(businessRegistry, "sales_total", "Total sales")
	
	// Infrastructure metrics - use with Collector middleware
	infraGauge := metrics.GaugeWithRegistry(infraRegistry, "cpu_usage", "CPU usage")

	// Use different registries in different routers/middleware
	// The metrics will be collected in their respective registries via Collector middleware
}
```

### API Reference

#### Custom Registry Functions

All metric creation functions have `*WithRegistry` variants:

- `CounterWithRegistry(registry *prometheus.Registry, name, help string) CounterMetric`
- `CounterWithRegistryWith[T any](registry *prometheus.Registry, name, help string) CounterMetricLabeled[T]`
- `GaugeWithRegistry(registry *prometheus.Registry, name, help string) GaugeMetric`
- `GaugeWithRegistryWith[T any](registry *prometheus.Registry, name, help string) GaugeMetricLabeled[T]`
- `HistogramWithRegistry(registry *prometheus.Registry, name, help string, buckets []float64) HistogramMetric`
- `HistogramWithRegistryWith[T any](registry *prometheus.Registry, name, help string, buckets []float64) HistogramMetricLabeled[T]`

#### Configuration Options

Both `CollectorOpts` and `TransportOpts` now support a `Registry` field:

```go
type CollectorOpts struct {
	Registry *prometheus.Registry // Optional custom registry
	Host     bool                 // Track host label
	Proto    bool                 // Track protocol label  
	Skip     func(*http.Request) bool // Skip function
}

type TransportOpts struct {
	Registry *prometheus.Registry // Optional custom registry
	Host     bool                 // Track host label (high cardinality warning!)
}
```

**Note**: If `Registry` is `nil`, the default global registry is used, ensuring full backwards compatibility.

## Example

See [_example/main.go](./_example/main.go) and try it locally:
```sh
$ cd _example

$ go run .
```

TODO: Run Prometheus + Grafana locally.

## License
[MIT license](./LICENSE)
