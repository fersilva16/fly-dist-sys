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

var node = maelstrom.NewNode()

func main() {
	messages := []int{}
	neighbors := []string{}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if !slices.Contains(messages, body.Message) {
			messages = append(messages, body.Message)

			for _, neighbor := range neighbors {
				neighborMsg := BroadcastRequest{
					MessageBody: maelstrom.MessageBody{
						Type: "broadcast",
					},

					Message: body.Message,
				}

				err := node.Send(neighbor, neighborMsg)

				if err != nil {
					return err
				}
			}
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
		var body TopologyRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		neighbors = body.Topology[node.ID()]

		resBody := map[string]any{
			"type": "topology_ok",
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
