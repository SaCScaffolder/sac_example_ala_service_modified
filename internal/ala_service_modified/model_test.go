package ala_service_modified

import (
	"strings"
	"testing"
)

// TestValidateAcceptsValidAlias confirms that a well-formed alias
// passes validation and gets lowercased + trimmed.
func TestValidateAcceptsValidAlias(t *testing.T) {
	a := &Alias{AliasURL: "  GitHub  ", RedirectURI: "https://github.com"}
	key, err := a.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "github" {
		t.Errorf("expected key=github, got %q", key)
	}
}

// TestValidateRejectsShortAlias enforces the 3-char minimum.
func TestValidateRejectsShortAlias(t *testing.T) {
	a := &Alias{AliasURL: "ab", RedirectURI: "https://example.com"}
	if _, err := a.Validate(); err != ErrInvalidAlias {
		t.Fatalf("expected ErrInvalidAlias, got %v", err)
	}
}

// TestValidateRejectsNonHttpRedirect enforces the http(s) scheme.
func TestValidateRejectsNonHttpRedirect(t *testing.T) {
	a := &Alias{AliasURL: "valid", RedirectURI: "ftp://example.com"}
	if _, err := a.Validate(); err != ErrInvalidRedirect {
		t.Fatalf("expected ErrInvalidRedirect, got %v", err)
	}
}

// TestNewShortCodeLength checks that the random generator produces
// 8-char base64url strings (6 bytes encoded without padding).
func TestNewShortCodeLength(t *testing.T) {
	for i := 0; i < 100; i++ {
		code := NewShortCode()
		if len(code) == 0 {
			t.Fatal("got empty short code (crypto/rand failure?)")
		}
		if len(code) != 8 {
			t.Errorf("expected length 8, got %d (code=%q)", len(code), code)
		}
		if strings.ContainsAny(code, "+/=") {
			t.Errorf("expected base64url (no padding), got %q", code)
		}
	}
}