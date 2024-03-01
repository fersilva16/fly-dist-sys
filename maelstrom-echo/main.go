package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type EchoRequest struct {
	maelstrom.MessageBody
	Echo string `json:"echo"`;
}

func main() {
  node := maelstrom.NewNode()

  node.Handle("echo", func (msg maelstrom.Message) error {
    var body EchoRequest;

    if err := json.Unmarshal(msg.Body, &body); err != nil {
      return err;
    }

    res_body := map[string]any{
      "type": "echo_ok",
      "echo": body.Echo,
    };

    return node.Reply(msg, res_body)
  });

  if err := node.Run(); err != nil {
    log.Fatal(err);
  }
}