package main

import (
	"log"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type GenerateResponse struct {
	maelstrom.MessageBody
	ID string `json:"id"`
}

var node = maelstrom.NewNode()
var clock Clock = TimeClock{}
var count int64

func main() {
	node.Handle("generate", func(msg maelstrom.Message) error {
		id := node.ID() + clock.Now() + strconv.FormatInt(count, 10)

		count += 1

		resBody := GenerateResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "generate_ok",
			},

			ID: id,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
