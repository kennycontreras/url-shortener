package main

import "sync"

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

func (s *URLStore) Get(key string) string {
	s.mu.RLock() // multiple go routines can read at a time but only one go routine can write at a time
	url := s.urls[key]
	s.mu.Unlock()
	return url
}

func (s *URLStore) Set(key, url string) bool {
	s.mu.Lock() // only one go routine can write at a time
	
	_, present := s.urls[key]
	if present {
		s.mu.Unlock()
		return false
	} else {
		s.urls[key] = url
		s.mu.Unlock()
		return true
	}
}

func main() {

}
