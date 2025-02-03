package main

import (
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

func (r *Replicator) Replicate(clock int, txnId int, key int, value interface{}) {
	for _, neighborId := range node.NodeIDs() {
		if neighborId == node.ID() {
			continue
		}

		nodeReplicator := r.load(neighborId)

		nodeReplicator.Replicate(clock, txnId, key, value)
	}
}

func (r *Replicator) load(neighborId string) *NodeReplicator {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeReplicator, ok := r.m[neighborId]

	if !ok {
		nodeReplicator = NewNodeReplicator(node, neighborId)

		r.m[neighborId] = nodeReplicator
	}

	return nodeReplicator
}

func (r *Replicator) Remove(src string, keys []int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeReplicator, ok := r.m[src]

	if !ok {
		return
	}

	nodeReplicator.Remove(keys)
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

func (r *NodeReplicator) Replicate(clock int, txnId int, key int, value interface{}) {
	r.m.Store(key, Value{
		Value: value,
		TxnId: txnId,
	})

	snapshot := map[int]Value{}

	r.m.Range(func(key, value any) bool {
		snapshot[key.(int)] = value.(Value)

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

func (r *NodeReplicator) Remove(keys []int) {
	for _, key := range keys {
		r.m.Delete(key)
	}
}
