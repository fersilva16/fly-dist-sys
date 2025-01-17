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
		txnStore := NewStore()

		for _, op := range body.Txn {
			if op.fn == READ {
				committedTxn = append(committedTxn, Op{READ, op.key, store.Read(op.key)})

				continue
			}

			txnStore.Write(txnId, op.key, op.value)

			committedTxn = append(committedTxn, op)
		}

		store.Merge(txnStore)

		clock := counter.GetLocal()

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

		counter.Sync(msg.Src, body.Clock)

		txnIds := []int{}
		txnStore := NewStore()

		for txnId, txn := range body.Snapshot {
			for _, op := range txn {
				txnStore.Write(txnId, op.key, op.value)
			}

			txnIds = append(txnIds, txnId)
		}

		store.Merge(txnStore)

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

		replicator.Remove(msg.Src, body.TxnIds)

		return nil
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
