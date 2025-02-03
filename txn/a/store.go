package main

import "sync"

type Store struct {
	m  map[int]interface{}
	mu sync.Mutex
}

func NewStore() *Store {
	return &Store{
		m: map[int]interface{}{},
	}
}

func (s *Store) Read(key int) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	value := s.m[key]

	return value
}

func (s *Store) Write(key int, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value
}
