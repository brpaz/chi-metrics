package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	// Example 1: Using the default registry (backwards compatible)
	fmt.Println("=== Example 1: Default Registry (Backwards Compatible) ===")
	defaultRegistryExample()

	// Example 2: Using a custom registry for all metrics
	fmt.Println("\n=== Example 2: Custom Registry for All Metrics ===")
	customRegistryExample()

	// Example 3: Using separate registries for different services
	fmt.Println("\n=== Example 3: Separate Registries for Different Services ===")
	separateRegistriesExample()
}

func defaultRegistryExample() {
	// This is the traditional way - all metrics go to the default registry
	r := chi.NewRouter()

	// Use the default collector (uses default registry)
	r.Use(metrics.Collector(metrics.CollectorOpts{
		Host:  false,
		Proto: true,
	}))

	// Use the default handler (serves from default registry)
	r.Handle("/metrics", metrics.Handler())

	// Create metrics using default registry
	counter := metrics.Counter("app_requests_total", "Total app requests")
	gauge := metrics.Gauge("app_active_connections", "Active connections")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		counter.Inc()
		gauge.Set(42)
		w.Write([]byte("Hello from default registry!\n"))
	})

	fmt.Println("Default registry server would run on :8080")
	fmt.Println("Metrics available at: http://localhost:8080/metrics")
}

func customRegistryExample() {
	// Create a custom registry
	customRegistry := prometheus.NewRegistry()

	r := chi.NewRouter()

	// Use collector with custom registry
	r.Use(metrics.Collector(metrics.CollectorOpts{
		Registry: customRegistry,
		Host:     true,
		Proto:    true,
	}))

	// Use handler from default registry (custom registry metrics collected via Collector)
	r.Handle("/metrics", metrics.Handler())

	// Create metrics in the custom registry
	counter := metrics.CounterWithRegistry(customRegistry, "custom_app_requests_total", "Total custom app requests")
	gauge := metrics.GaugeWithRegistry(customRegistry, "custom_app_active_connections", "Custom active connections")

	// Create typed labeled metrics in custom registry
	type requestLabels struct {
		Method string `label:"method"`
		Status string `label:"status"`
	}
	labeledCounter := metrics.CounterWithRegistryWith[requestLabels](
		customRegistry,
		"custom_app_requests_by_method_total",
		"Total requests by method and status",
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		counter.Inc()
		gauge.Set(100)
		labeledCounter.Inc(requestLabels{Method: "GET", Status: "200"})
		w.Write([]byte("Hello from custom registry!\n"))
	})

	fmt.Println("Custom registry server would run on :8081")
	fmt.Println("Custom metrics available at: http://localhost:8081/metrics")
}

func separateRegistriesExample() {
	// Create separate registries for different concerns
	appRegistry := prometheus.NewRegistry()
	infraRegistry := prometheus.NewRegistry()

	// App server with app-specific metrics
	appRouter := chi.NewRouter()
	appRouter.Use(metrics.Collector(metrics.CollectorOpts{
		Registry: appRegistry,
		Host:     false,
		Proto:    true,
	}))
	appRouter.Handle("/metrics", metrics.Handler())

	// Infrastructure server with infra-specific metrics
	infraRouter := chi.NewRouter()
	infraRouter.Use(metrics.Collector(metrics.CollectorOpts{
		Registry: infraRegistry,
		Host:     true,
		Proto:    false,
	}))
	infraRouter.Handle("/metrics", metrics.Handler())

	// Create app-specific metrics
	appCounter := metrics.CounterWithRegistry(appRegistry, "business_transactions_total", "Business transactions")
	appGauge := metrics.GaugeWithRegistry(appRegistry, "business_active_sessions", "Active user sessions")

	// Create infrastructure metrics
	infraCounter := metrics.CounterWithRegistry(infraRegistry, "system_calls_total", "System calls")
	infraGauge := metrics.GaugeWithRegistry(infraRegistry, "system_memory_usage", "Memory usage")

	// Custom transport for outgoing requests with separate registry
	transport := metrics.Transport(metrics.TransportOpts{
		Registry: infraRegistry,
		Host:     true,
	})
	client := &http.Client{
		Transport: transport(http.DefaultTransport),
	}

	appRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		appCounter.Inc()
		appGauge.Set(250)

		// Make an outgoing request (tracked in infra registry)
		go func() {
			resp, err := client.Get("http://example.com")
			if err == nil {
				resp.Body.Close()
			}
		}()

		w.Write([]byte("Business logic completed!\n"))
	})

	infraRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		infraCounter.Inc()
		infraGauge.Set(1024)
		w.Write([]byte("Infrastructure healthy!\n"))
	})

	fmt.Println("App server would run on :8082 - business metrics")
	fmt.Println("Infra server would run on :8083 - infrastructure metrics")
	fmt.Println("App metrics: http://localhost:8082/metrics")
	fmt.Println("Infra metrics: http://localhost:8083/metrics")
}

// Example of migration from default to custom registry
func migrationExample() {
	fmt.Println("\n=== Migration Example ===")

	// Step 1: Existing code (no changes needed)
	oldCounter := metrics.Counter("old_metric", "Old metric")
	oldCounter.Inc()

	// Step 2: Gradually introduce custom registry
	customRegistry := prometheus.NewRegistry()
	newCounter := metrics.CounterWithRegistry(customRegistry, "new_metric", "New metric in custom registry")
	newCounter.Inc()

	// Step 3: Registries are isolated via Collector middleware
	// Only the default handler is available
	defaultHandler := metrics.Handler()              // serves default registry

	fmt.Printf("Default handler type: %T\n", defaultHandler)
	fmt.Println("Custom registries work independently via Collector middleware!")
}

// Demonstrates nil registry behavior (falls back to default)
func nilRegistryExample() {
	fmt.Println("\n=== Nil Registry Example (Default Fallback) ===")

	// Passing nil registry is equivalent to using the original functions
	counter1 := metrics.Counter("example_counter", "Example counter")
	counter2 := metrics.CounterWithRegistry(nil, "example_counter_2", "Example counter 2")

	handler1 := metrics.Handler()

	fmt.Printf("Counter1 type: %T\n", counter1)
	fmt.Printf("Counter2 type: %T\n", counter2)
	fmt.Printf("Handler1 type: %T\n", handler1)
	fmt.Println("Nil registry functions are equivalent to original functions!")
}