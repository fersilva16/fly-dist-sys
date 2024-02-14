package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	echo(n);

	uniqueIds(n);

	broadcast(n);

	if err := n.Run(); err != nil {
		log.Fatal(err);
	}
}