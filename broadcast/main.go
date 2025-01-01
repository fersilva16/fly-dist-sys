package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastRequest struct {
	maelstrom.MessageBody
	Message int `json:"message"`
}

type TopologyRequest struct {
	maelstrom.MessageBody
	Topology map[string][]string `json:"topology"`
}

var node = maelstrom.NewNode()

func main() {
	messages := []int{}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages = append(messages, body.Message)

		resBody := map[string]any{
			"type": "broadcast_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		resBody := map[string]any{
			"type":     "read_ok",
			"messages": messages,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		resBody := map[string]any{
			"type": "topology_ok",
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
