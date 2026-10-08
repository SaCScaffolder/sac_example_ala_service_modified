package ala_service_modified

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCORSHandlerSetsHeaders asserts the middleware sets the standard
// CORS response headers on a regular request.
func TestCORSHandlerSetsHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	CORSHandler(next).ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected AllowOrigin=*, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Errorf("expected AllowMethods to be set")
	}
}

// TestCORSHandlerShortCircuitsOptions checks that an OPTIONS request
// returns 204 without calling the wrapped handler.
func TestCORSHandlerShortCircuitsOptions(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/v1/aliases", nil)

	CORSHandler(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204 on OPTIONS, got %d", rec.Code)
	}
	if called {
		t.Error("expected wrapped handler NOT to be called on OPTIONS")
	}
}