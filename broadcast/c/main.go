package main

import (
	"encoding/json"

	"log"
	"slices"
	"time"

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

			for _, neighborId := range neighbors {
				if neighborId == msg.Src {
					continue
				}

				go broadcast(500*time.Millisecond, neighborId, body.Message)
			}
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

func broadcast(timeout time.Duration, neighborId string, message int) {
	replied := false

	neighborMessage := BroadcastRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "broadcast",
		},

		Message: message,
	}

	node.RPC(neighborId, neighborMessage, func(msg maelstrom.Message) error {
		replied = true

		return nil
	})

	time.Sleep(timeout)

	if !replied {
		broadcast(timeout*2, neighborId, message)
	}
}
