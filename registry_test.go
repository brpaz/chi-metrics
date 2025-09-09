package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestCustomRegistryHandler(t *testing.T) {
	registry := prometheus.NewRegistry()

	// Create a metric in the custom registry
	counter := CounterWithRegistry(registry, "test_handler_counter", "Test counter for handler")
	counter.Inc()

	// Create handler with custom registry
	handler := HandlerWithRegistry(registry)

	// Test the handler
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "test_handler_counter") {
		t.Errorf("Expected response to contain 'test_handler_counter', but it didn't. Body: %s", body)
	}

	if !strings.Contains(body, "1") {
		t.Errorf("Expected response to contain counter value '1', but it didn't. Body: %s", body)
	}
}

func TestCustomRegistryHandlerWithNil(t *testing.T) {
	// Test that HandlerWithRegistry(nil) behaves like Handler()
	defaultHandler := Handler()
	nilRegistryHandler := HandlerWithRegistry(nil)

	// Both should work (though we can't easily test they're identical)
	req := httptest.NewRequest("GET", "/metrics", nil)
	
	w1 := httptest.NewRecorder()
	defaultHandler.ServeHTTP(w1, req)
	
	w2 := httptest.NewRecorder()
	nilRegistryHandler.ServeHTTP(w2, req)

	// Both should return 200 OK
	if w1.Result().StatusCode != http.StatusOK {
		t.Errorf("Default handler returned status %d", w1.Result().StatusCode)
	}
	if w2.Result().StatusCode != http.StatusOK {
		t.Errorf("Nil registry handler returned status %d", w2.Result().StatusCode)
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