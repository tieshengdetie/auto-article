package utils

import "sync"

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func (s *SafeMap) Get(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[key]
}

func (s *SafeMap) Set(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}
