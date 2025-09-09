package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestCustomRegistryCounter(t *testing.T) {
	registry := prometheus.NewRegistry()

	// Test Counter with custom registry
	counter := CounterWithRegistry(registry, "test_counter", "Test counter")
	counter.Inc()
	counter.Add(5.0)

	// Verify the metric is registered in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metricFamilies) != 1 {
		t.Fatalf("Expected 1 metric family, got %d", len(metricFamilies))
	}

	if *metricFamilies[0].Name != "test_counter" {
		t.Errorf("Expected metric name 'test_counter', got '%s'", *metricFamilies[0].Name)
	}

	if *metricFamilies[0].Metric[0].Counter.Value != 6.0 {
		t.Errorf("Expected counter value 6.0, got %f", *metricFamilies[0].Metric[0].Counter.Value)
	}
}

func TestCustomRegistryCounterWithLabels(t *testing.T) {
	registry := prometheus.NewRegistry()

	type testLabels struct {
		Method string `label:"method"`
		Status string `label:"status"`
	}

	// Test CounterWith with custom registry
	counter := CounterWithRegistryWith[testLabels](registry, "test_counter_labeled", "Test counter with labels")
	
	labels := testLabels{Method: "GET", Status: "200"}
	counter.Inc(labels)
	counter.Add(3.0, labels)

	// Verify the metric is registered in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metricFamilies) != 1 {
		t.Fatalf("Expected 1 metric family, got %d", len(metricFamilies))
	}

	if *metricFamilies[0].Name != "test_counter_labeled" {
		t.Errorf("Expected metric name 'test_counter_labeled', got '%s'", *metricFamilies[0].Name)
	}

	if *metricFamilies[0].Metric[0].Counter.Value != 4.0 {
		t.Errorf("Expected counter value 4.0, got %f", *metricFamilies[0].Metric[0].Counter.Value)
	}

	// Check labels
	metric := metricFamilies[0].Metric[0]
	if len(metric.Label) != 2 {
		t.Fatalf("Expected 2 labels, got %d", len(metric.Label))
	}

	labelMap := make(map[string]string)
	for _, label := range metric.Label {
		labelMap[*label.Name] = *label.Value
	}

	if labelMap["method"] != "GET" {
		t.Errorf("Expected method label 'GET', got '%s'", labelMap["method"])
	}
	if labelMap["status"] != "200" {
		t.Errorf("Expected status label '200', got '%s'", labelMap["status"])
	}
}

func TestCustomRegistryGauge(t *testing.T) {
	registry := prometheus.NewRegistry()

	// Test Gauge with custom registry
	gauge := GaugeWithRegistry(registry, "test_gauge", "Test gauge")
	gauge.Set(10.5)
	gauge.Inc()
	gauge.Dec()
	gauge.Add(2.5)

	// Verify the metric is registered in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metricFamilies) != 1 {
		t.Fatalf("Expected 1 metric family, got %d", len(metricFamilies))
	}

	if *metricFamilies[0].Name != "test_gauge" {
		t.Errorf("Expected metric name 'test_gauge', got '%s'", *metricFamilies[0].Name)
	}

	if *metricFamilies[0].Metric[0].Gauge.Value != 13.0 {
		t.Errorf("Expected gauge value 13.0, got %f", *metricFamilies[0].Metric[0].Gauge.Value)
	}
}

func TestCustomRegistryGaugeWithLabels(t *testing.T) {
	registry := prometheus.NewRegistry()

	type testLabels struct {
		Service string `label:"service"`
	}

	// Test GaugeWith with custom registry
	gauge := GaugeWithRegistryWith[testLabels](registry, "test_gauge_labeled", "Test gauge with labels")
	
	labels := testLabels{Service: "api"}
	gauge.Set(5.0, labels)
	gauge.Inc(labels)

	// Verify the metric is registered in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metricFamilies) != 1 {
		t.Fatalf("Expected 1 metric family, got %d", len(metricFamilies))
	}

	if *metricFamilies[0].Metric[0].Gauge.Value != 6.0 {
		t.Errorf("Expected gauge value 6.0, got %f", *metricFamilies[0].Metric[0].Gauge.Value)
	}
}

func TestCustomRegistryHistogram(t *testing.T) {
	registry := prometheus.NewRegistry()

	// Test Histogram with custom registry
	buckets := []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0}
	histogram := HistogramWithRegistry(registry, "test_histogram", "Test histogram", buckets)
	
	histogram.Observe(0.5)
	histogram.Observe(1.5)
	histogram.Observe(3.0)

	// Verify the metric is registered in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metricFamilies) != 1 {
		t.Fatalf("Expected 1 metric family, got %d", len(metricFamilies))
	}

	if *metricFamilies[0].Name != "test_histogram" {
		t.Errorf("Expected metric name 'test_histogram', got '%s'", *metricFamilies[0].Name)
	}

	if *metricFamilies[0].Metric[0].Histogram.SampleCount != 3 {
		t.Errorf("Expected histogram sample count 3, got %d", *metricFamilies[0].Metric[0].Histogram.SampleCount)
	}
}

func TestCustomRegistryHistogramWithLabels(t *testing.T) {
	registry := prometheus.NewRegistry()

	type testLabels struct {
		Method string `label:"method"`
	}

	// Test HistogramWith with custom registry
	buckets := []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0}
	histogram := HistogramWithRegistryWith[testLabels](registry, "test_histogram_labeled", "Test histogram with labels", buckets)
	
	labels := testLabels{Method: "POST"}
	histogram.Observe(1.2, labels)
	histogram.Observe(0.8, labels)

	// Verify the metric is registered in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metricFamilies) != 1 {
		t.Fatalf("Expected 1 metric family, got %d", len(metricFamilies))
	}

	if *metricFamilies[0].Metric[0].Histogram.SampleCount != 2 {
		t.Errorf("Expected histogram sample count 2, got %d", *metricFamilies[0].Metric[0].Histogram.SampleCount)
	}
}



func TestNilRegistryFallback(t *testing.T) {
	// Test that nil registry falls back to default registry behavior
	
	// These should not panic and should work like the original functions
	counter := CounterWithRegistry(nil, "test_nil_counter", "Test nil registry counter")
	counter.Inc()

	gauge := GaugeWithRegistry(nil, "test_nil_gauge", "Test nil registry gauge")
	gauge.Set(42.0)

	buckets := []float64{0.1, 1.0, 10.0}
	histogram := HistogramWithRegistry(nil, "test_nil_histogram", "Test nil registry histogram", buckets)
	histogram.Observe(5.0)

	type testLabels struct {
		Label string `label:"label"`
	}

	counterLabeled := CounterWithRegistryWith[testLabels](nil, "test_nil_counter_labeled", "Test nil registry counter labeled")
	counterLabeled.Inc(testLabels{Label: "test"})

	gaugeLabeled := GaugeWithRegistryWith[testLabels](nil, "test_nil_gauge_labeled", "Test nil registry gauge labeled")
	gaugeLabeled.Set(123.0, testLabels{Label: "test"})

	histogramLabeled := HistogramWithRegistryWith[testLabels](nil, "test_nil_histogram_labeled", "Test nil registry histogram labeled", buckets)
	histogramLabeled.Observe(2.5, testLabels{Label: "test"})

	// If we reach here without panicking, the test passes
	t.Log("All nil registry fallback tests passed")
}

func TestBackwardsCompatibility(t *testing.T) {
	// Test that the original functions still work exactly as before
	
	counter := Counter("backwards_counter", "Backwards compatibility counter")
	counter.Inc()
	counter.Add(10.0)

	gauge := Gauge("backwards_gauge", "Backwards compatibility gauge")
	gauge.Set(100.0)
	gauge.Inc()
	gauge.Dec()

	buckets := []float64{0.1, 1.0, 10.0}
	histogram := Histogram("backwards_histogram", "Backwards compatibility histogram", buckets)
	histogram.Observe(1.5)

	type testLabels struct {
		Type string `label:"type"`
	}

	counterLabeled := CounterWith[testLabels]("backwards_counter_labeled", "Backwards compatibility counter labeled")
	counterLabeled.Inc(testLabels{Type: "test"})

	gaugeLabeled := GaugeWith[testLabels]("backwards_gauge_labeled", "Backwards compatibility gauge labeled")
	gaugeLabeled.Set(50.0, testLabels{Type: "test"})

	histogramLabeled := HistogramWith[testLabels]("backwards_histogram_labeled", "Backwards compatibility histogram labeled", buckets)
	histogramLabeled.Observe(0.5, testLabels{Type: "test"})

	// Test Handler function
	handler := Handler()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Backwards compatibility handler returned status %d", w.Result().StatusCode)
	}

	t.Log("All backwards compatibility tests passed")
}

func TestCustomRegistryCollector(t *testing.T) {
	registry := prometheus.NewRegistry()
	
	// Create a collector with custom registry
	collector := Collector(CollectorOpts{
		Registry: registry,
		Host:     true,
		Proto:    true,
	})

	// Create a test handler
	handler := collector(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Check metrics are in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Should have metrics: requests_total, request_duration_seconds, and possibly inflight
	if len(metricFamilies) < 2 {
		t.Errorf("Expected at least 2 metric families, got %d", len(metricFamilies))
	}

	metricNames := make(map[string]bool)
	for _, mf := range metricFamilies {
		metricNames[*mf.Name] = true
		t.Logf("Found metric: %s", *mf.Name)
	}

	if !metricNames["http_requests_total"] {
		t.Error("Expected http_requests_total metric")
	}
	if !metricNames["http_request_duration_seconds"] {
		t.Error("Expected http_request_duration_seconds metric")
	}
}

func TestCustomRegistryTransport(t *testing.T) {
	registry := prometheus.NewRegistry()
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server response"))
	}))
	defer server.Close()

	// Create transport with custom registry
	transport := Transport(TransportOpts{
		Registry: registry,
		Host:     true,
	})

	// Create client with custom transport
	client := &http.Client{
		Transport: transport(http.DefaultTransport),
	}

	// Make a request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	resp.Body.Close()

	// Check metrics are in the custom registry
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Should have metrics: client_requests_total, client_request_duration_seconds, and possibly inflight
	if len(metricFamilies) < 2 {
		t.Errorf("Expected at least 2 metric families, got %d", len(metricFamilies))
	}

	metricNames := make(map[string]bool)
	for _, mf := range metricFamilies {
		metricNames[*mf.Name] = true
		t.Logf("Found metric: %s", *mf.Name)
	}

	if !metricNames["http_client_requests_total"] {
		t.Error("Expected http_client_requests_total metric")
	}
	if !metricNames["http_client_request_duration_seconds"] {
		t.Error("Expected http_client_request_duration_seconds metric")
	}
}

func TestCollectorWithNilRegistry(t *testing.T) {
	// Test that collector with nil registry works like the original
	collector := Collector(CollectorOpts{
		Registry: nil,
		Host:     false,
		Proto:    false,
	})

	// Create a test handler
	handler := collector(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should not panic and should complete successfully
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Result().StatusCode)
	}
}

func TestTransportWithNilRegistry(t *testing.T) {
	// Test that transport with nil registry works like the original
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server response"))
	}))
	defer server.Close()

	// Create transport with nil registry
	transport := Transport(TransportOpts{
		Registry: nil,
		Host:     false,
	})

	// Create client with custom transport
	client := &http.Client{
		Transport: transport(http.DefaultTransport),
	}

	// Make a request - should not panic
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}