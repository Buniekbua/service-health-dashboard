package storage

import (
	"maps"
	"sync"
)

type Storage struct {
	mu       sync.Mutex
	statuses map[string]int
}

func NewStorage() *Storage {
	return &Storage{
		statuses: make(map[string]int),
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
