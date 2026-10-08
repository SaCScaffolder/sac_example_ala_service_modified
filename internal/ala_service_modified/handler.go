package ala_service_modified

import (
	"encoding/json"
	"net/http"
	"strings"
)

type createAliasRequest struct {
	AliasURL    string `json:"alias_url"`
	RedirectURI string `json:"redirect_uri"`
}

type aliasResponse struct {
	AliasURL    string `json:"alias_url"`
	RedirectURI string `json:"redirect_uri"`
}

// CreateAlias handles POST /v1/aliases. Validates the request,
// stores the alias, and returns 201 with the canonical form.
func (s *Server) CreateAlias(w http.ResponseWriter, r *http.Request) {
	var req createAliasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	a := &Alias{AliasURL: req.AliasURL, RedirectURI: req.RedirectURI}
	key, err := a.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.AliasURL = key
	if err := s.store.Put(a); err != nil {
		if err == ErrAliasExists {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, aliasResponse{
		AliasURL:    a.AliasURL,
		RedirectURI: a.RedirectURI,
	})
}

// GetAliases handles GET /v1/aliases and returns the full list.
func (s *Server) GetAliases(w http.ResponseWriter, r *http.Request) {
	all := s.store.List()
	writeJSON(w, http.StatusOK, map[string]any{
		"count":  len(all),
		"values": all,
	})
}

// Redirect handles GET /<alias_url>. Looks up the alias and returns
// 302 with Location, or 404 if not found.
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		s.index(w, r)
		return
	}
	a, ok := s.store.Get(key)
	if !ok {
		writeError(w, http.StatusNotFound, ErrAliasNotFound.Error())
		return
	}
	http.Redirect(w, r, a.RedirectURI, http.StatusFound)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}