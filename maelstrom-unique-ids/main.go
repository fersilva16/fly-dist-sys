package main

import (
	"log"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
  node := maelstrom.NewNode()

  var count int64;

  node.Handle("generate", func(msg maelstrom.Message) error {  
		id := node.ID() + strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(count, 10);

    count += 1;

    res_body := map[string]any{
      "type": "generate_ok",
      "id": id,
    };
    
    return node.Reply(msg, res_body);
  })
  
  if err := node.Run(); err != nil {
    log.Fatal(err);
  }
}