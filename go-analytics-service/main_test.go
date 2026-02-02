package main

import (
	"testing"
	"time"
)

func TestRollingAverage(t *testing.T) {
	service, err := NewService("localhost:6379")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Test with known values
	values := []float64{100, 110, 105, 115, 120, 95, 100, 105, 110, 115}
	expectedAvg := 107.5

	for _, v := range values {
		service.processMetricValue(v)
	}

	service.analytics.mu.RLock()
	prediction := service.analytics.prediction
	service.analytics.mu.RUnlock()

	tolerance := 0.1
	if prediction < expectedAvg-tolerance || prediction > expectedAvg+tolerance {
		t.Errorf("Expected prediction ~%.2f, got %.2f", expectedAvg, prediction)
	}
}

func TestZScoreAnomalyDetection(t *testing.T) {
	service, err := NewService("localhost:6379")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Fill window with normal values
	for i := 0; i < windowSize; i++ {
		service.processMetricValue(100.0)
	}

	initialAnomalies := service.analytics.anomalyCount

	// Send anomaly (high value)
	service.processMetricValue(300.0) // Should be detected as anomaly

	service.analytics.mu.RLock()
	anomalyCount := service.analytics.anomalyCount
	service.analytics.mu.RUnlock()

	if anomalyCount <= initialAnomalies {
		t.Error("Anomaly should have been detected")
	}
}

func TestMetricStructure(t *testing.T) {
	metric := Metric{
		Timestamp: time.Now(),
		CPU:       50.5,
		RPS:       100.0,
	}

	if metric.CPU != 50.5 {
		t.Errorf("Expected CPU 50.5, got %.2f", metric.CPU)
	}

	if metric.RPS != 100.0 {
		t.Errorf("Expected RPS 100.0, got %.2f", metric.RPS)
	}
}

