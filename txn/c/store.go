package main

import (
	"sync"
)

type Value struct {
	Value interface{} `json:"value"`
	TxnId int         `json:"txnId"`
}

type Store struct {
	m  map[int]*Value
	mu sync.Mutex
}

func NewStore() *Store {
	return &Store{
		m: make(map[int]*Value),
	}
}

func (s *Store) Read(key int) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.m[key]

	if !ok {
		return nil
	}

	return value.Value
}

func (s *Store) Write(txnId int, key int, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.write(txnId, key, value)
}

func (s *Store) write(txnId int, key int, value interface{}) {
	currentValue, ok := s.m[key]

	if ok && currentValue.TxnId > txnId {
		return
	}

	s.m[key] = &Value{
		Value: value,
		TxnId: txnId,
	}
}

func (s *Store) Range(fn func(key int, value *Value) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range s.m {
		if !fn(k, v) {
			break
		}
	}
}

func (s *Store) Merge(other *Store) {
	s.mu.Lock()
	defer s.mu.Unlock()

	other.Range(func(key int, value *Value) bool {
		s.write(value.TxnId, key, value.Value)

		return true
	})
}
