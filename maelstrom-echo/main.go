package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
  node := maelstrom.NewNode()

  node.Handle("echo", func (msg maelstrom.Message) error {
    var body map[string]any;

    if err := json.Unmarshal(msg.Body, &body); err != nil {
      return err;
    }

    res_body := map[string]any{
      "type": "echo_ok",
      "echo": body["echo"],
    };

    return node.Reply(msg, res_body)
  });

  if err := node.Run(); err != nil {
    log.Fatal(err);
  }
}