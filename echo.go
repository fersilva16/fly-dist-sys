package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func echo(n *maelstrom.Node) {
	n.Handle("echo", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		body["type"] = "echo_ok";

		return n.Reply(msg, body)
	})
}