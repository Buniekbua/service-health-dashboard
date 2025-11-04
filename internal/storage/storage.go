package storage

import (
	"maps"
	"sync"
)

type Storage struct {
	mu       sync.Mutex
	statuses map[string]int
	urls     map[string]struct{}
}

func NewStorage() *Storage {
	return &Storage{
		statuses: make(map[string]int),
		urls:     make(map[string]struct{}),
	}
}

func (s *Storage) UpdateStatus(url string, status int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statuses[url] = status
}

func (s *Storage) GetAllStatuses() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make(map[string]int, len(s.statuses))
	maps.Copy(copy, s.statuses)
	return copy
}

func (s *Storage) AddURL(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[url] = struct{}{}
}

func (s *Storage) RemoveURL(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.urls, url)
	delete(s.statuses, url)
}

func (s *Storage) GetURLs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	urls := make([]string, 0, len(s.urls))
	for u := range s.urls {
		urls = append(urls, u)
	}
	return urls
}
