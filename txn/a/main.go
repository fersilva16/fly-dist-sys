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

var node = maelstrom.NewNode()

func main() {
	store := NewStore()

	node.Handle("txn", func(msg maelstrom.Message) error {
		var body TxnRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		committedTxn := Txn{}

		for _, op := range body.Txn {
			if op.fn == READ {
				committedTxn = append(committedTxn, Op{READ, op.key, store.Read(op.key)})

				continue
			}

			store.Write(op.key, op.value)

			committedTxn = append(committedTxn, op)
		}

		resBody := TxnResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "txn_ok",
			},

			Txn: committedTxn,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
