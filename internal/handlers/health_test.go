package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/umar/helloworld/internal/handlers"
)

func TestHealthz(t *testing.T) {
	t.Run("returns 200 ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()

		handlers.Healthz(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var body map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["status"] != "ok" {
			t.Errorf("expected status ok, got %q", body["status"])
		}
	})

	t.Run("POST method still returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
		rec := httptest.NewRecorder()

		handlers.Healthz(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for POST, got %d", rec.Code)
		}
	})
}

func TestReadyz(t *testing.T) {
	t.Run("returns 200 ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()

		handlers.Readyz(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var body map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["status"] != "ok" {
			t.Errorf("expected status ok, got %q", body["status"])
		}
	})

	t.Run("content-type is application/json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()

		handlers.Readyz(rec, req)

		ct := rec.Header().Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
	})
}
