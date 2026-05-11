package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/umar/helloworld/internal/handlers"
)

func newTestHelloHandler(t *testing.T) *handlers.HelloHandler {
	t.Helper()
	reg := prometheus.NewRegistry()
	return handlers.NewHelloHandler(reg)
}

func TestHelloHandler_WithSecret(t *testing.T) {
	t.Setenv("HELLO_SECRET", "supersecret")

	h := newTestHelloHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["message"] != "supersecret" {
		t.Errorf("expected message supersecret, got %q", body["message"])
	}
}

func TestHelloHandler_EmptySecret(t *testing.T) {
	t.Setenv("HELLO_SECRET", "")

	h := newTestHelloHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 even when secret is empty, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if _, ok := body["message"]; !ok {
		t.Error("expected message key in response body")
	}
}

func TestHelloHandler_ContentType(t *testing.T) {
	h := newTestHelloHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHandleFail(t *testing.T) {
	h := newTestHelloHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	rec := httptest.NewRecorder()

	h.HandleFail(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["error"] == "" {
		t.Error("expected error key in response body")
	}
}
