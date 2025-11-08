package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupMockServer creates a mock HTTP server that returns JSON or an error code
func setupMockServer(status int, body any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetSchedule_Success(t *testing.T) {
	mockData := []ScheduleResponse{MockScheduleResponse()}
	server := setupMockServer(http.StatusOK, mockData)
	defer server.Close()

	client := &http.Client{}
	internalId := "user123"

	got, err := GetSchedule(client, server.URL, internalId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(got))
	}
}

func TestGetSchedule_MissingInternalID(t *testing.T) {
	client := &http.Client{}
	_, err := GetSchedule(client, "http://example.com", "")
	if err == nil || err.Error() != "internalId is required" {
		t.Errorf("expected internalId error, got %v", err)
	}
}

func TestGetSchedule_HTTPError(t *testing.T) {
	server := setupMockServer(http.StatusInternalServerError, nil)
	defer server.Close()

	client := &http.Client{}
	_, err := GetSchedule(client, server.URL, "abc")
	if err == nil {
		t.Error("expected error for HTTP 500, got nil")
	}
}

func TestGetSchedule_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := GetSchedule(client, server.URL, "abc")
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
