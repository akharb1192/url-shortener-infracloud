package store

import "testing"

func TestShortenIdempotent(t *testing.T) {
	s := NewInMemoryStore()
	url := "https://example.com/path"
	c1, err := s.Shorten(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c2, err := s.Shorten(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c1 != c2 {
		t.Fatalf("expected same code for same url")
	}
}

func TestResolve(t *testing.T) {
	s := NewInMemoryStore()
	url := "https://example.com/hello"
	code, _ := s.Shorten(url)
	got, err := s.Resolve(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != url {
		t.Fatalf("expected %s got %s", url, got)
	}
}

func TestTopDomains(t *testing.T) {
	s := NewInMemoryStore()
	s.Shorten("https://youtube.com/1")
	s.Shorten("https://youtube.com/2")
	s.Shorten("https://udemy.com/1")
	s.Shorten("https://udemy.com/2")
	s.Shorten("https://udemy.com/3")

	list := s.TopDomains(3)
	if list[0].Domain != "udemy.com" || list[0].Count != 3 {
		t.Fatalf("unexpected top domain %+v", list[0])
	}
}
