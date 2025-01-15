package main

import "sync"

type Store struct {
	m  map[int]interface{}
	mu sync.Mutex
}

func NewStore() *Store {
	return &Store{
		m: make(map[int]interface{}),
	}
}

func (s *Store) Read(key int) (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.m[key]

	return value, ok
}

func (s *Store) Write(key int, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value
}
