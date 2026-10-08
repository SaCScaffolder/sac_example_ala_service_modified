package ala_service_modified

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateAliasHappyPath asserts a valid POST creates a 201 response
// and stores the alias in the store.
func TestCreateAliasHappyPath(t *testing.T) {
	srv := NewServer(NewMemoryStore())
	body := `{"alias_url":"github","redirect_uri":"https://github.com"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/aliases", strings.NewReader(body))
	rec := httptest.NewRecorder()

	srv.v1Aliases(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", rec.Code, rec.Body.String())
	}

	var got aliasResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.AliasURL != "github" {
		t.Errorf("expected alias_url=github, got %q", got.AliasURL)
	}
	if got.RedirectURI != "https://github.com" {
		t.Errorf("expected redirect_uri=https://github.com, got %q", got.RedirectURI)
	}
}

// TestCreateAliasRejectsMissingUrl checks the validation path catches
// bad input before it reaches the store.
func TestCreateAliasRejectsMissingUrl(t *testing.T) {
	srv := NewServer(NewMemoryStore())
	body := `{"alias_url":"","redirect_uri":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/aliases", strings.NewReader(body))
	rec := httptest.NewRecorder()

	srv.v1Aliases(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

// TestCreateAliasDuplicateReturns409 covers the conflict path when
// the same alias_url is registered twice.
func TestCreateAliasDuplicateReturns409(t *testing.T) {
	store := NewMemoryStore()
	srv := NewServer(store)

	// First request: 201.
	body := `{"alias_url":"xxx","redirect_uri":"https://example.com"}`
	req1 := httptest.NewRequest(http.MethodPost, "/v1/aliases", strings.NewReader(body))
	srv.v1Aliases(httptest.NewRecorder(), req1)

	// Second request: same alias_url, different redirect_uri.
	body2 := `{"alias_url":"xxx","redirect_uri":"https://other.com"}`
	req2 := httptest.NewRequest(http.MethodPost, "/v1/aliases", strings.NewReader(body2))
	rec := httptest.NewRecorder()
	srv.v1Aliases(rec, req2)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 on duplicate, got %d (body=%s)", rec.Code, rec.Body.String())
	}
}

// TestRedirectHappyPath asserts GET /<alias> returns 302 with Location.
func TestRedirectHappyPath(t *testing.T) {
	store := NewMemoryStore()
	if err := store.Put(&Alias{AliasURL: "gh", RedirectURI: "https://github.com"}); err != nil {
		t.Fatalf("seed store: %v", err)
	}
	srv := NewServer(store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/gh", nil)

	srv.routeRoot(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://github.com" {
		t.Errorf("expected Location=https://github.com, got %q", loc)
	}
}

// TestRedirectNotFound asserts a missing alias returns 404.
func TestRedirectNotFound(t *testing.T) {
	srv := NewServer(NewMemoryStore())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/nope", nil)

	srv.routeRoot(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}