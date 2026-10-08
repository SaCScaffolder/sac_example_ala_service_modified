package ala_service_modified

import (
	"sync"
	"testing"
)

// TestMemoryStorePutGet is the smoke test: store an alias, retrieve it.
func TestMemoryStorePutGet(t *testing.T) {
	s := NewMemoryStore()
	a := &Alias{AliasURL: "x", RedirectURI: "https://example.com"}
	if err := s.Put(a); err != nil {
		t.Fatalf("put: %v", err)
	}
	got, ok := s.Get("x")
	if !ok {
		t.Fatal("expected ok=true after Put")
	}
	if got.RedirectURI != "https://example.com" {
		t.Errorf("expected redirect_uri=https://example.com, got %q", got.RedirectURI)
	}
}

// TestMemoryStorePutDuplicateReturnsErr checks that re-inserting the
// same alias_url fails with ErrAliasExists.
func TestMemoryStorePutDuplicateReturnsErr(t *testing.T) {
	s := NewMemoryStore()
	if err := s.Put(&Alias{AliasURL: "x", RedirectURI: "https://a.com"}); err != nil {
		t.Fatalf("first put: %v", err)
	}
	if err := s.Put(&Alias{AliasURL: "x", RedirectURI: "https://b.com"}); err != ErrAliasExists {
		t.Fatalf("expected ErrAliasExists, got %v", err)
	}
}

// TestMemoryStoreListReturnsAll confirms List returns every alias.
func TestMemoryStoreListReturnsAll(t *testing.T) {
	s := NewMemoryStore()
	for _, k := range []string{"a", "b", "c"} {
		if err := s.Put(&Alias{AliasURL: k, RedirectURI: "https://" + k + ".com"}); err != nil {
			t.Fatalf("put %s: %v", k, err)
		}
	}
	all := s.List()
	if len(all) != 3 {
		t.Errorf("expected 3 aliases, got %d", len(all))
	}
}

// TestMemoryStoreConcurrentSafe runs concurrent Put+Get+List to verify
// the store is safe for use under -race.
func TestMemoryStoreConcurrentSafe(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "k"
			if i%2 == 0 {
				key = "k"
			} else {
				key = "j"
			}
			_ = s.Put(&Alias{AliasURL: key, RedirectURI: "https://x.com"})
			_, _ = s.Get(key)
			_ = s.List()
		}(i)
	}
	wg.Wait()
}