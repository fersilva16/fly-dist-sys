package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type EchoRequest struct {
	maelstrom.MessageBody
	Echo string `json:"echo"`
}

type EchoResponse struct {
	maelstrom.MessageBody
	Echo string `json:"echo"`
}

var node = maelstrom.NewNode()

func main() {
	node.Handle("echo", func(msg maelstrom.Message) error {
		var body EchoRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resBody := EchoResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "echo_ok",
			},

			Echo: body.Echo,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
