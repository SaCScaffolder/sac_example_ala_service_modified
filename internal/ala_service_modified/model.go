package ala_service_modified

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

// Alias represents a short-URL -> long-URL mapping.
type Alias struct {
	AliasURL    string    `json:"alias_url"`
	RedirectURI string    `json:"redirect_uri"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate enforces invariants on an alias before it's stored.
// Returns the lowercased, trimmed alias_url ready for insertion.
func (a *Alias) Validate() (string, error) {
	url := strings.ToLower(strings.TrimSpace(a.AliasURL))
	if len(url) < 3 || len(url) > 64 {
		return "", ErrInvalidAlias
	}
	if !strings.HasPrefix(a.RedirectURI, "http://") && !strings.HasPrefix(a.RedirectURI, "https://") {
		return "", ErrInvalidRedirect
	}
	a.AliasURL = url
	return url, nil
}

// NewShortCode returns a base64url-encoded 6-byte random string.
// Used by tests to verify generator invariants.
func NewShortCode() string {
	var buf [6]byte
	if _, err := rand.Read(buf[:]); err != nil {
		// crypto/rand failure on Linux is essentially impossible;
		// a fixed fallback would be a security regression. Return
		// an empty string and let the caller handle it.
		return ""
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(buf[:]), "=")
}

var (
	ErrInvalidAlias    = errorString("alias_url must be 3-64 chars")
	ErrInvalidRedirect = errorString("redirect_uri must be http(s)")
	ErrAliasExists     = errorString("alias_url already exists")
	ErrAliasNotFound   = errorString("alias_url not found")
)

type errorString string

func (e errorString) Error() string { return string(e) }