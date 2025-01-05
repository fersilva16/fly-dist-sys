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
	neighbors := []string{}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body types.BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if !slices.Contains(messages, body.Message) {
			messages = append(messages, body.Message)

			for _, neighborId := range neighbors {
				if neighborId == msg.Src {
					continue
				}

				neighborMessage := types.BroadcastRequest{
					MessageBody: maelstrom.MessageBody{
						Type: "broadcast",
					},

					Message: body.Message,
				}

				node.Send(neighborId, neighborMessage)
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
		var body types.TopologyRequest

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
