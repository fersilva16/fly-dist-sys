package main

import (
	"encoding/json"
	"log"
	"slices"
	"sync"

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
var neighbours = map[string]*sync.Map{}

func main() {
	var messages []int

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if slices.Contains(messages, body.Message) {
			return nil
		}

		messages = append(messages, body.Message)

		for id, messages := range neighbours {
			neighbourMessage := map[string]any{
				"type":    "broadcast",
				"message": body.Message,
			}

			messages.Store(body.MsgID, neighbourMessage)

			go func(messages *sync.Map, id string) {
				messages.Range(func(key any, rawValue any) bool {
					message := rawValue.(map[string]any)

					node.RPC(id, message, func(replyMsg maelstrom.Message) error {
						messages.Delete(key)

						return nil
					})

					return true
				})
			}(messages, id)
		}

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
		var body TopologyRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology := body.Topology[node.ID()]

		for i := 0; i < len(topology); i++ {
			id := topology[i]

			neighbours[id] = &sync.Map{}
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
