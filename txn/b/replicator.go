package main

import (
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Replicator struct {
	node *maelstrom.Node
	m    sync.Map
}

func NewReplicator(node *maelstrom.Node) *Replicator {
	return &Replicator{
		node: node,
		m:    sync.Map{},
	}
}

func (r *Replicator) Replicate(clock int, txnId int, txn Txn) {
	r.add(txnId, txn)

	for _, neighbor := range node.NodeIDs() {
		if neighbor == node.ID() {
			continue
		}

		snapshot := map[int]Txn{}

		r.m.Range(func(key, value any) bool {
			snapshot[key.(int)] = value.(Txn)

			return true
		})

		neighborMessage := ReplicateRequest{
			MessageBody: maelstrom.MessageBody{
				Type: "replicate",
			},

			Clock:    clock,
			Snapshot: snapshot,
		}

		go func() {
			node.Send(neighbor, neighborMessage)
		}()
	}
}

func (r *Replicator) add(txnId int, txn Txn) {
	// Only replicate writes
	writeOnlyTxn := Txn{}

	for _, op := range txn {
		if op.fn == WRITE {
			writeOnlyTxn = append(writeOnlyTxn, op)
		}
	}

	r.m.Store(txnId, writeOnlyTxn)
}

func (r *Replicator) Remove(txnIds []int) {
	for _, txnId := range txnIds {
		r.m.Delete(txnId)
	}
}
