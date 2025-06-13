package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"analytics-dashboard-api/internal/handlers"
)

func TestHealthHandler_Health(t *testing.T) {
	logger := &mockLogger{}
	handler := handlers.NewHealthHandler(logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.Health(recorder, req)

	// Check status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Health() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	// Check content type
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Health() content-type = %s, want application/json", contentType)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Health() response parsing error: %v", err)
	}

	// Check required fields
	requiredFields := []string{"status", "timestamp", "uptime", "version", "memory", "goroutines"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Health() missing required field: %s", field)
		}
	}

	// Check status field
	if status, ok := response["status"].(string); !ok || status != "healthy" {
		t.Errorf("Health() status = %v, want 'healthy'", response["status"])
	}

	// Check version field
	if version, ok := response["version"].(string); !ok || version != "1.0.0" {
		t.Errorf("Health() version = %v, want '1.0.0'", response["version"])
	}

	// Check memory structure
	if memory, ok := response["memory"].(map[string]interface{}); ok {
		memoryFields := []string{"alloc_mb", "total_alloc_mb", "sys_mb", "num_gc"}
		for _, field := range memoryFields {
			if _, exists := memory[field]; !exists {
				t.Errorf("Health() memory missing field: %s", field)
			}
		}
	} else {
		t.Error("Health() memory field should be an object")
	}

	// Check goroutines field
	if goroutines, ok := response["goroutines"].(float64); !ok || goroutines < 0 {
		t.Errorf("Health() goroutines should be a positive number, got %v", response["goroutines"])
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	logger := &mockLogger{}
	handler := handlers.NewHealthHandler(logger)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	recorder := httptest.NewRecorder()

	handler.Ready(recorder, req)

	// Check status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Ready() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	// Check content type
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Ready() content-type = %s, want application/json", contentType)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Ready() response parsing error: %v", err)
	}

	// Check required fields
	requiredFields := []string{"status", "timestamp"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Ready() missing required field: %s", field)
		}
	}

	// Check status field
	if status, ok := response["status"].(string); !ok || status != "ready" {
		t.Errorf("Ready() status = %v, want 'ready'", response["status"])
	}

	// Check timestamp format (should be parseable as RFC3339)
	if timestampStr, ok := response["timestamp"].(string); ok {
		if _, err := time.Parse(time.RFC3339, timestampStr); err != nil {
			t.Errorf("Ready() timestamp format invalid: %s", timestampStr)
		}
	} else {
		t.Error("Ready() timestamp should be a string")
	}
}

func TestHealthHandler_HealthUptime(t *testing.T) {
	logger := &mockLogger{}
	handler := handlers.NewHealthHandler(logger)

	// Wait a small amount to ensure uptime is positive
	time.Sleep(1 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.Health(recorder, req)

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Health() response parsing error: %v", err)
	}

	// Check uptime is a valid duration string
	if uptimeStr, ok := response["uptime"].(string); ok {
		if _, err := time.ParseDuration(uptimeStr); err != nil {
			t.Errorf("Health() uptime format invalid: %s", uptimeStr)
		}
	} else {
		t.Error("Health() uptime should be a string")
	}
}
