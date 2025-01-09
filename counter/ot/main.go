package main

import (
	"context"
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type AddRequest struct {
	maelstrom.MessageBody
	Delta int `json:"delta"`
}

var node = maelstrom.NewNode()

func main() {
	kv := maelstrom.NewSeqKV(node)

	node.Handle("add", func(msg maelstrom.Message) error {
		var body AddRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()

		for {
			value, err := kv.ReadInt(ctx, "counter")

			if err != nil {
				value = 0
			}

			err = kv.CompareAndSwap(ctx, "counter", value, value+body.Delta, true)

			if err == nil {
				resBody := map[string]any{
					"type": "add_ok",
				}

				return node.Reply(msg, resBody)
			}
		}
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		ctx := context.Background()

		value, err := kv.ReadInt(ctx, "counter")

		if err != nil {
			value = 0
		}

		resBody := map[string]any{
			"type":  "read_ok",
			"value": value,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
