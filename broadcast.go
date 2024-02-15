package main

import (
	"encoding/json"
	"slices"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func broadcast(n *maelstrom.Node) {
	var messages []any;
	neighbours := map[string]*sync.Map{};

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any;

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
		}

		if (slices.Contains(messages, body["message"])) {
			return nil;
		}

		messages = append(messages, body["message"]);
		
		for id, messages := range neighbours {			
			neighbour_message := map[string]any{
				"type": "broadcast",
				"message": body["message"],
			}
			
			messages.Store(body["msg_id"], neighbour_message);

			go func(messages *sync.Map, id string) {
				messages.Range(func (key any, raw_value any) bool {
					message := raw_value.(map[string]any);
					
					n.RPC(id, message, func (reply_msg maelstrom.Message) error {
						messages.Delete(key);
						
						return nil;
					});
					
					return true;
				});
			}(messages, id);
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
			id := topology[i].(string)

			neighbours[id] = &sync.Map{};
		}

		res_body := map[string]any{
			"type": "topology_ok",
		}

		return n.Reply(msg, res_body);
	})
}