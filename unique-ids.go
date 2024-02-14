package main

import (
	"encoding/json"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func uniqueIds(n *maelstrom.Node) {
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
}