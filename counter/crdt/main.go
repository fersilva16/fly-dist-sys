package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type AddRequest struct {
	maelstrom.MessageBody
	Delta int `json:"delta"`
}

type PropagateRequest struct {
	maelstrom.MessageBody
	Count int `json:"count"`
}

type ReadResponse struct {
	maelstrom.MessageBody
	Value int `json:"value"`
}

var node = maelstrom.NewNode()

func main() {
	crdt := NewCRDT()
	count := 0

	node.Handle("add", func(msg maelstrom.Message) error {
		var body AddRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		count += body.Delta

		for _, neighborId := range node.NodeIDs() {
			if neighborId == node.ID() {
				continue
			}

			neighborMessage := PropagateRequest{
				MessageBody: maelstrom.MessageBody{
					Type: "propagate",
				},

				Count: count,
			}

			go func() {
				node.Send(neighborId, neighborMessage)
			}()
		}

		resBody := maelstrom.MessageBody{
			Type: "add_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("propagate", func(msg maelstrom.Message) error {
		var body PropagateRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		crdt.Sync(msg.Src, body.Count)

		return nil
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		value := count + crdt.Read()

		resBody := ReadResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "read_ok",
			},

			Value: value,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
