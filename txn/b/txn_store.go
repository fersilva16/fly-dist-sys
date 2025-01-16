package main

import (
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Value struct {
	value interface{}
	txnId int
}

type TxnStore struct {
	m       map[int]*Value
	mu      sync.Mutex
	counter *GCounter
}

func NewTxnStore(node *maelstrom.Node) *TxnStore {
	return &TxnStore{
		m:       make(map[int]*Value),
		counter: NewGCounter(node),
	}
}

func (s *TxnStore) Merge(src string, clock int, snapshot map[int]Txn) []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter.Sync(src, clock)

	committedTxnIds := []int{}

	for tnxId, txn := range snapshot {
		for _, op := range txn {
			s.commitOp(tnxId, op)
		}

		committedTxnIds = append(committedTxnIds, tnxId)
	}

	return committedTxnIds
}

func (s *TxnStore) Commit(txn Txn) (int, Txn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	txnId := s.counter.Increment()

	committedTxn := Txn{}

	for _, op := range txn {
		committedTxn = append(committedTxn, s.commitOp(txnId, op))
	}

	return txnId, committedTxn
}

func (s *TxnStore) commitOp(txnId int, op Op) Op {
	if op.fn == WRITE {
		return s.commitWrite(txnId, op.key, op.value)
	}

	return s.commitRead(op.key)
}

func (s *TxnStore) commitWrite(txnId, key int, value interface{}) Op {
	currentValue, ok := s.m[key]

	if ok && currentValue.txnId > txnId {
		return Op{READ, key, currentValue.value}
	}

	s.m[key] = &Value{
		value: value,
		txnId: txnId,
	}

	return Op{WRITE, key, value}
}

func (s *TxnStore) commitRead(key int) Op {
	value, ok := s.m[key]

	if !ok {
		return Op{READ, key, nil}
	}

	return Op{READ, key, value.value}
}
