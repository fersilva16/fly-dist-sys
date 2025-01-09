package main

import (
	"encoding/json"

	"log"
	"slices"

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

type ReadResponse struct {
	maelstrom.MessageBody
	Messages []int `json:"messages"`
}

var node = maelstrom.NewNode()

func main() {
	messages := []int{}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if !slices.Contains(messages, body.Message) {
			messages = append(messages, body.Message)
		}

		resBody := maelstrom.MessageBody{
			Type: "broadcast_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		return nil
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		resBody := ReadResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "read_ok",
			},

			Messages: messages,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body TopologyRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resBody := maelstrom.MessageBody{
			Type: "topology_ok",
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
