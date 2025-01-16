package main

import "sync"

type TxnStore struct {
	m  map[int]interface{}
	mu sync.Mutex
}

func NewTxnStore() *TxnStore {
	return &TxnStore{
		m: make(map[int]interface{}),
	}
}

func (s *TxnStore) Commit(txn Txn) Txn {
	s.mu.Lock()
	defer s.mu.Unlock()

	committedTxn := Txn{}

	for _, op := range txn {
		committedTxn = append(committedTxn, s.commitOp(op))
	}

	return committedTxn
}

func (s *TxnStore) commitOp(op Op) Op {
	if op.fn == WRITE {
		return s.commitWrite(op.key, op.value)
	}

	return s.commitRead(op.key)
}

func (s *TxnStore) commitWrite(key int, value interface{}) Op {
	s.m[key] = value

	return Op{WRITE, key, value}
}

func (s *TxnStore) commitRead(key int) Op {
	value := s.m[key]

	return Op{READ, key, value}
}
