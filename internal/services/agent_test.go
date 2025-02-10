package services

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

// TestMetricAgentService_Start_SendingMetrics tests that metrics are collected and sent correctly
func TestMetricAgentService_Start_SendingMetrics(t *testing.T) {
	var receivedRequests []string

	// Create a test HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRequests = append(receivedRequests, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Create the MetricAgentService with the test server's URL
	service := &MetricAgentService{
		MetricChannel:  make(chan types.UpdateMetricValueRequest, 100),
		PollInterval:   1 * time.Second,
		ReportInterval: 1 * time.Second,
		Shutdown:       make(chan os.Signal, 1),
		APIClient:      resty.New(),
		Address:        testServer.URL,
	}

	// Add a test metric to the channel
	service.MetricChannel <- types.UpdateMetricValueRequest{
		Type:  types.Gauge,
		Name:  "Alloc",
		Value: "12345",
	}
	service.MetricChannel <- types.UpdateMetricValueRequest{
		Type:  types.Gauge,
		Name:  "BuckHashSys",
		Value: "1444382",
	}
	service.MetricChannel <- types.UpdateMetricValueRequest{
		Type:  types.Gauge,
		Name:  "Frees",
		Value: "127",
	}

	// Test sending metrics
	go service.Start()

	// Allow time for the metric to be sent
	time.Sleep(2 * time.Second)

	// Assert that the expected metric URL is present in the received requests
	expectedURLs := []string{
		"/update/gauge/Alloc/12345",
		"/update/gauge/BuckHashSys/1444382",
		"/update/gauge/Frees/127",
	}

	for _, expectedURL := range expectedURLs {
		require.Contains(t, receivedRequests, expectedURL, "Expected URL not found: %s", expectedURL)
	}
}

// TestMetricAgentService_Start_ErrorHandling tests error handling when sending metrics (e.g., 500 HTTP response)
func TestMetricAgentService_Start_ErrorHandling(t *testing.T) {
	var receivedRequests []string

	// Create a test HTTP server that returns a 500 error
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRequests = append(receivedRequests, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError) // Respond with a 500 error
	}))
	defer testServer.Close()

	// Create the MetricAgentService with the test server's URL
	service := &MetricAgentService{
		MetricChannel:  make(chan types.UpdateMetricValueRequest, 100),
		PollInterval:   1 * time.Second,
		ReportInterval: 1 * time.Second,
		Shutdown:       make(chan os.Signal, 1),
		APIClient:      resty.New(),
		Address:        testServer.URL,
	}

	// Add a test metric to the channel
	service.MetricChannel <- types.UpdateMetricValueRequest{
		Type:  types.Gauge,
		Name:  "Alloc",
		Value: "12345",
	}

	// Test sending metrics
	go service.Start()

	// Allow time for the metric to be sent
	time.Sleep(2 * time.Second)

	// Assert that the expected metric URL is present in the received requests
	require.Contains(t, receivedRequests, "/update/gauge/Alloc/12345")

	// Check that the error handling was triggered by checking the logs
	// (This might require capturing the logs if you're using a logger)
}

// TestMetricAgentService_Shutdown tests the graceful shutdown of the service
func TestMetricAgentService_Shutdown(t *testing.T) {
	// Create a test HTTP server that returns a status code 200
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a status OK (200)
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Create the MetricAgentService with the test server's URL
	service := &MetricAgentService{
		MetricChannel:  make(chan types.UpdateMetricValueRequest, 100),
		PollInterval:   1 * time.Second,
		ReportInterval: 1 * time.Second,
		Shutdown:       make(chan os.Signal, 1),
		APIClient:      resty.New(),
		Address:        testServer.URL,
	}

	// Start the service in a goroutine to simulate metric collection
	go service.Start()

	// Simulate shutdown
	service.Shutdown <- os.Interrupt

	// Allow time for shutdown processing
	time.Sleep(1 * time.Second)

	// Assert that the service shuts down gracefully (no panic or hanging)
	// This could be checked by ensuring no further metrics are being sent
}
