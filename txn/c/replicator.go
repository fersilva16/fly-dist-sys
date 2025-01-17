package main

import (
	"fmt"
	"os"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Replicator struct {
	node *maelstrom.Node
	m    map[string]*NodeReplicator
	mu   sync.Mutex
}

func NewReplicator(node *maelstrom.Node) *Replicator {
	return &Replicator{
		node: node,
		m:    map[string]*NodeReplicator{},
	}
}

func (r *Replicator) Replicate(clock int, txnId int, txn Txn) {
	writeOnlyTxn := Txn{}

	for _, op := range txn {
		if op.fn == WRITE {
			writeOnlyTxn = append(writeOnlyTxn, op)
		}
	}

	if len(writeOnlyTxn) == 0 {
		return
	}

	for _, neighbor := range node.NodeIDs() {
		if neighbor == node.ID() {
			continue
		}

		nodeReplicator := r.load(neighbor)

		nodeReplicator.Replicate(clock, txnId, writeOnlyTxn)
	}
}

func (r *Replicator) load(neighbor string) *NodeReplicator {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeReplicator, ok := r.m[neighbor]

	if !ok {
		nodeReplicator = NewNodeReplicator(node, neighbor)

		r.m[neighbor] = nodeReplicator
	}

	return nodeReplicator
}

func (r *Replicator) Remove(src string, txnIds []int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeReplicator, ok := r.m[src]

	if !ok {
		return
	}

	nodeReplicator.Remove(txnIds)
}

type NodeReplicator struct {
	node *maelstrom.Node
	id   string
	m    sync.Map
}

func NewNodeReplicator(node *maelstrom.Node, id string) *NodeReplicator {
	return &NodeReplicator{
		node: node,
		id:   id,
		m:    sync.Map{},
	}
}

func (r *NodeReplicator) Replicate(clock int, txnId int, txn Txn) {
	r.m.Store(txnId, txn)

	snapshot := map[int]Txn{}

	r.m.Range(func(txnId, value any) bool {
		snapshot[txnId.(int)] = value.(Txn)

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
		node.Send(r.id, neighborMessage)
	}()
}

func (r *NodeReplicator) Remove(txnIds []int) {
	for _, txnId := range txnIds {
		fmt.Fprintf(os.Stderr, "Removing txn %d from %s\n", txnId, r.id)
		r.m.Delete(txnId)
	}
}
