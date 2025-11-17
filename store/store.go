package store

import (
	"errors"
	"net/url"
	"strings"
	"sync"
)

type InMemoryStore struct {
	mu           sync.RWMutex
	u2c          map[string]string
	c2u          map[string]string
	domainCounts map[string]int
	counter      uint64
}

var (
	ErrNotFound = errors.New("not found")
)

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		u2c:          make(map[string]string),
		c2u:          make(map[string]string),
		domainCounts: make(map[string]int),
	}
}

func (s *InMemoryStore) Shorten(rawURL string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", errors.New("invalid url")
	}
	norm := u.String()

	s.mu.RLock()
	if code, ok := s.u2c[norm]; ok {
		s.mu.RUnlock()
		s.incrementDomainCount(u.Hostname())
		return code, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if code, ok := s.u2c[norm]; ok {
		s.incrementDomainCount(u.Hostname())
		return code, nil
	}

	s.counter++
	code := encodeBase62(s.counter)
	s.u2c[norm] = code
	s.c2u[code] = norm
	s.domainCounts[u.Hostname()]++

	return code, nil
}

func (s *InMemoryStore) incrementDomainCount(domain string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.domainCounts[domain]++
}

func (s *InMemoryStore) Resolve(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, ok := s.c2u[code]; ok {
		return u, nil
	}
	return "", ErrNotFound
}

type DomainCount struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

func (s *InMemoryStore) TopDomains(n int) []DomainCount {
	s.mu.RLock()
	m := make(map[string]int, len(s.domainCounts))
	for k, v := range s.domainCounts {
		m[k] = v
	}
	s.mu.RUnlock()

	list := make([]DomainCount, 0, len(m))
	for k, v := range m {
		list = append(list, DomainCount{Domain: k, Count: v})
	}

	sortDomainCounts(list)

	if n > len(list) {
		n = len(list)
	}
	return list[:n]
}

func sortDomainCounts(list []DomainCount) {
	for i := 0; i < len(list); i++ {
		maxIdx := i
		for j := i + 1; j < len(list); j++ {
			if list[j].Count > list[maxIdx].Count {
				maxIdx = j
			}
		}
		list[i], list[maxIdx] = list[maxIdx], list[i]
	}
}
