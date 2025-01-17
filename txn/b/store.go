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

func (s *Store) Write(txnId int, key int, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentValue, ok := s.m[key]

	if ok && currentValue.TxnId > txnId {
		return
	}

	s.m[key] = &Value{
		Value: value,
		TxnId: txnId,
	}
}
