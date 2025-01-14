package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var node = maelstrom.NewNode()

func main() {
	// TODO

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
