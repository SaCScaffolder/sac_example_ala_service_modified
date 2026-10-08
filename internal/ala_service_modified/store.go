package ala_service_modified

import "sync"

// Store is the alias persistence layer. The in-memory implementation
// here is sufficient for the integration-test fixture; production
// services should swap in a sql.DB-backed implementation.
type Store interface {
	Put(a *Alias) error
	Get(aliasURL string) (Alias, bool)
	List() []Alias
}

// NewMemoryStore returns a thread-safe in-memory Store keyed by
// lowercased alias_url.
func NewMemoryStore() Store {
	return &memoryStore{aliases: map[string]Alias{}}
}

type memoryStore struct {
	mu      sync.RWMutex
	aliases map[string]Alias
}

func (m *memoryStore) Put(a *Alias) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.aliases[a.AliasURL]; ok {
		return ErrAliasExists
	}
	m.aliases[a.AliasURL] = *a
	return nil
}

func (m *memoryStore) Get(aliasURL string) (Alias, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.aliases[aliasURL]
	return a, ok
}

func (m *memoryStore) List() []Alias {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Alias, 0, len(m.aliases))
	for _, a := range m.aliases {
		out = append(out, a)
	}
	return out
}