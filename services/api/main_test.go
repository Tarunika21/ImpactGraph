package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandlerReturnsHealthy(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected HTTP 200, got HTTP %d",
			recorder.Code,
		)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Error("expected JSON content type")
	}

	var response HealthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("expected healthy, got %q", response.Status)
	}

	if response.Service != "impactgraph-api" {
		t.Errorf("unexpected service name: %q", response.Service)
	}

	if _, err := time.Parse(time.RFC3339, response.Timestamp); err != nil {
		t.Errorf("invalid timestamp: %v", err)
	}
}

func TestHealthHandlerRejectsPost(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected HTTP 405, got HTTP %d",
			recorder.Code,
		)
	}

	if recorder.Header().Get("Allow") != http.MethodGet {
		t.Error("expected Allow header GET")
	}
}
