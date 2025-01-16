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

	Clock    int         `json:"clock"`
	Snapshot map[int]Txn `json:"snapshot"`
}

type ReplicateResponse struct {
	maelstrom.MessageBody

	TxnIds []int `json:"txnIds"`
}

var node = maelstrom.NewNode()

func main() {
	store := NewTxnStore(node)
	replicator := NewReplicator(node)

	node.Handle("txn", func(msg maelstrom.Message) error {
		var body TxnRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		txnId, committedTxn := store.Commit(body.Txn)

		clock := store.counter.GetLocal()

		replicator.Replicate(clock, txnId, committedTxn)

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

		txnIds := store.Merge(msg.Src, body.Clock, body.Snapshot)

		resBody := ReplicateResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "replicate_ok",
			},

			TxnIds: txnIds,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("replicate_ok", func(msg maelstrom.Message) error {
		var body ReplicateResponse

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		replicator.Remove(body.TxnIds)

		return nil
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
