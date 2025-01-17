package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TxnRequest struct {
	maelstrom.MessageBody

	Txn Txn `json:"txn"`
}

type TxnResponse struct {
	maelstrom.MessageBody

	Txn Txn `json:"txn"`
}

type ReplicateRequest struct {
	maelstrom.MessageBody

	Clock    int           `json:"clock"`
	Snapshot map[int]Value `json:"snapshot"`
}

type ReplicateResponse struct {
	maelstrom.MessageBody

	Keys []int `json:"keys"`
}

var node = maelstrom.NewNode()

func main() {
	counter := NewGCounter(node)
	store := NewStore()
	replicator := NewReplicator(node)

	node.Handle("txn", func(msg maelstrom.Message) error {
		var body TxnRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		committedTxn := Txn{}

		txnId := counter.Increment()

		for _, op := range body.Txn {
			if op.fn == READ {
				committedTxn = append(committedTxn, Op{READ, op.key, store.Read(op.key)})

				continue
			}

			store.Write(txnId, op.key, op.value)

			committedTxn = append(committedTxn, op)

			clock := counter.GetLocal()

			replicator.Replicate(clock, txnId, op.key, op.value)
		}

		resBody := TxnResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "txn_ok",
			},

			Txn: committedTxn,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("replicate", func(msg maelstrom.Message) error {
		var body ReplicateRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		counter.Sync(msg.Src, body.Clock)

		keys := []int{}

		for key, value := range body.Snapshot {
			store.Write(value.TxnId, key, value.Value)

			keys = append(keys, key)
		}

		resBody := ReplicateResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "replicate_ok",
			},

			Keys: keys,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("replicate_ok", func(msg maelstrom.Message) error {
		var body ReplicateResponse

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		replicator.Remove(body.Keys)

		return nil
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
