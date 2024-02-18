package main

import (
	"encoding/json"
	"log"
	"slices"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
  node := maelstrom.NewNode()

  var messages []any;
  neighbours := map[string]*sync.Map{};

  node.Handle("broadcast", func(msg maelstrom.Message) error {
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
          
          node.RPC(id, message, func (reply_msg maelstrom.Message) error {
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
    
    return node.Reply(msg, res_body);
  })

  node.Handle("read", func(msg maelstrom.Message) error {
    var body map[string]any;

    if err := json.Unmarshal(msg.Body, &body); err != nil {
      return err;
    }

    res_body := map[string]any{
      "type": "read_ok",
      "messages": messages,
    }

    return node.Reply(msg, res_body);
  })
  
  node.Handle("topology", func(msg maelstrom.Message) error {
    var body map[string]any;

    if err := json.Unmarshal(msg.Body, &body); err != nil {
      return err;
    }

    topology := body["topology"].(map[string]interface {})[node.ID()].([]interface{});

    for i := 0; i < len(topology); i++ {
      id := topology[i].(string)

      neighbours[id] = &sync.Map{};
    }

    res_body := map[string]any{
      "type": "topology_ok",
    }

    return node.Reply(msg, res_body);
  })
  
  if err := node.Run(); err != nil {
    log.Fatal(err);
  }
}