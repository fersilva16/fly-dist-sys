package main

import (
	"encoding/json"
	"gossip-gloomers/utils"

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

type GossipRequest struct {
	maelstrom.MessageBody
	Messages []int `json:"messages"`
}

type ReadResponse struct {
	maelstrom.MessageBody
	Messages []int `json:"messages"`
}

var node = maelstrom.NewNode()

func main() {
	set := utils.NewSet[int]()

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resBody := maelstrom.MessageBody{
			Type: "broadcast_ok",
		}

		set.Add(body.Message)

		messages := set.GetAll()

		for _, neighborId := range node.NodeIDs() {
			if neighborId == node.ID() {
				continue
			}

			neighborMessage := GossipRequest{
				MessageBody: maelstrom.MessageBody{
					Type: "gossip",
				},

				Messages: messages,
			}

			go func() {
				node.Send(neighborId, neighborMessage)
			}()
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("gossip", func(msg maelstrom.Message) error {
		var body GossipRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		set.AddAll(body.Messages)

		return nil
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		msgs := set.GetAll()

		resBody := ReadResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "read_ok",
			},

			Messages: msgs,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		resBody := maelstrom.MessageBody{
			Type: "topology_ok",
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
