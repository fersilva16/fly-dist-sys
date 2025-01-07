package utils

import "sync"

type Set[T comparable] struct {
	m  map[T]bool
	mu sync.Mutex
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: map[T]bool{},
	}
}

func (s *Set[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[item] = true
}

func (s *Set[T]) AddAll(items []T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		s.m[item] = true
	}
}

func (s *Set[T]) GetAll() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := []T{}

	for item := range s.m {
		items = append(items, item)
	}

	return items
}
