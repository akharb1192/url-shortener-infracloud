package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/akharb1192/url-shortener-infracloud/store"
)

type Handler struct {
	store *store.InMemoryStore
	base  string
}

func NewHandler(s *store.InMemoryStore) *Handler {
	return &Handler{store: s, base: "http://localhost:8080"}
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "invalid json", http.StatusBadRequest)
		return
	}

	code, err := h.store.Shorten(req.URL)
	if err != nil {
		h.respondError(w, "invalid url", http.StatusBadRequest)
		return
	}

	short := fmt.Sprintf("%s/s/%s", strings.TrimRight(h.base, "/"), code)
	h.respondJSON(w, map[string]string{"short_url": short})
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := path.Base(r.URL.Path)
	if code == "s" || code == "" {
		h.respondError(w, "missing code", http.StatusBadRequest)
		return
	}

	u, err := h.store.Resolve(code)
	if err != nil {
		h.respondError(w, "not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, u, http.StatusFound)
}

func (h *Handler) TopDomainsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list := h.store.TopDomains(3)
	h.respondJSON(w, list)
}

func (h *Handler) respondJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) respondError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	h.respondJSON(w, map[string]string{"error": msg})
}
