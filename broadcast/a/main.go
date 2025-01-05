package main

import (
	"encoding/json"
	types "fly-dist-sys/broadcast"
	"log"
	"slices"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var node = maelstrom.NewNode()

func main() {
	messages := []int{}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body types.BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if !slices.Contains(messages, body.Message) {
			messages = append(messages, body.Message)
		}

		resBody := map[string]any{
			"type": "broadcast_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		return nil
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		resBody := map[string]any{
			"type":     "read_ok",
			"messages": messages,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body types.TopologyRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resBody := map[string]any{
			"type": "topology_ok",
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
