package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("echo", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		body["type"] = "echo_ok";

		return n.Reply(msg, body)
	})

	var count int64;

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}
			
		body["type"] = "generate_ok";
		body["id"] = n.ID() + strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(count, 10);

		count += 1;
		
		return n.Reply(msg, body);
	})

	var messages []any;

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		messages = append(messages, body["message"])

		res_body := make(map[string]any);

		res_body["type"] = "broadcast_ok";
		
		return n.Reply(msg, res_body);
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		body["type"] = "read_ok";
		body["messages"] = messages;

		return n.Reply(msg, body);
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		body := make(map[string]any)

		body["type"] = "topology_ok";

		return n.Reply(msg, body);
	})
	if err := n.Run(); err != nil {
		log.Fatal(err);
	}
}