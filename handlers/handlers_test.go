package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akharb1192/url-shortener-infracloud/store"
)

func TestShortenAndRedirect(t *testing.T) {
	s := store.NewInMemoryStore()
	h := NewHandler(s)

	body := map[string]string{"url": "https://example.com/x"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.ShortenHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	short := resp["short_url"]

	parts := bytes.Split([]byte(short), []byte("/"))
	code := string(parts[len(parts)-1])

	req2 := httptest.NewRequest(http.MethodGet, "/s/"+code, nil)
	w2 := httptest.NewRecorder()
	h.RedirectHandler(w2, req2)

	if w2.Code != http.StatusFound {
		t.Fatalf("expected 302 got %d", w2.Code)
	}
}
