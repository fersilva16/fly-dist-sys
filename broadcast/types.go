package types

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"

type BroadcastRequest struct {
	maelstrom.MessageBody
	Message int `json:"message"`
}

type TopologyRequest struct {
	maelstrom.MessageBody
	Topology map[string][]string `json:"topology"`
}
