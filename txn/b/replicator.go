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

func (r *Replicator) Replicate(clock int, txnId int, key int, value interface{}) {
	r.m.Store(key, Value{
		Value: value,
		TxnId: txnId,
	})

	for _, neighbor := range node.NodeIDs() {
		if neighbor == node.ID() {
			continue
		}

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
			node.Send(neighbor, neighborMessage)
		}()
	}
}

func (r *Replicator) Remove(keys []int) {
	for _, key := range keys {
		r.m.Delete(key)
	}
}
