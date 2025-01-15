package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TxnRequest struct {
	maelstrom.MessageBody

	Txn []Op `json:"txn"`
}

type TxnResponse struct {
	maelstrom.MessageBody

	Txn []Op `json:"txn"`
}

var node = maelstrom.NewNode()

func main() {
	store := NewStore()

	node.Handle("txn", func(msg maelstrom.Message) error {
		var body TxnRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		res := []Op{}

		for _, op := range body.Txn {
			if op.fn == WRITE {
				store.Write(op.key, op.value)

				res = append(res, Op{WRITE, op.key, op.value})

				continue
			}

			value, ok := store.Read(op.key)

			if !ok {
				value = 0
			}

			res = append(res, Op{READ, op.key, value})
		}

		resBody := TxnResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "txn_ok",
			},

			Txn: res,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
