package main

import (
	"encoding/json"
	"slices"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func broadcast(n *maelstrom.Node) {
	var messages []any;
	var handled_messages []any;
	var neighbours []string;

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		if (slices.Contains(handled_messages, body["msg_id"])) {
			return nil;
		}

		handled_messages = append(handled_messages, body["msg_id"]);

		messages = append(messages, body["message"]);
		
		for i := 0; i < len(neighbours); i++ {
			n.Send(neighbours[i], map[string]any{
				"type": "broadcast",
				"message": body["message"],
				"msg_id": body["msg_id"],
			});
		}

		res_body := map[string]any{
			"type": "broadcast_ok",
		};
		
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
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		topology := body["topology"].(map[string]interface {})[n.ID()].([]interface{});

		for i := 0; i < len(topology); i++ {
			neighbours = append(neighbours, topology[i].(string))	
		}

		res_body := map[string]any{
			"type": "topology_ok",
		}

		return n.Reply(msg, res_body);
	})
}