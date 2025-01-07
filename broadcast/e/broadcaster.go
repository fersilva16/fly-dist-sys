package main

import (
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Broadcaster struct {
	node    *maelstrom.Node
	buffer  []int
	stopChn chan bool
	chn     chan int
	ticker  *time.Ticker
}

func NewBroadcaster(node *maelstrom.Node, d time.Duration) *Broadcaster {
	return &Broadcaster{
		node:    node,
		buffer:  []int{},
		chn:     make(chan int, 30),
		stopChn: make(chan bool, 1),
		ticker:  time.NewTicker(d),
	}
}

func (b *Broadcaster) Start() {
	for {
		select {
		case <-b.stopChn:
			b.ticker.Stop()
			close(b.chn)
			close(b.stopChn)
			return
		case <-b.ticker.C:
			b.broadcast()
		}
	}
}

func (b *Broadcaster) Append(message int) {
	b.chn <- message
}

func (b *Broadcaster) flush() bool {
	if len(b.chn) == 0 {
		return false
	}

	for len(b.chn) > 0 {
		b.buffer = append(b.buffer, <-b.chn)
	}

	return true
}

func (b *Broadcaster) broadcast() {
	b.flush()

	for _, neighbor := range node.NodeIDs() {
		if neighbor == node.ID() {
			continue
		}

		neighborMessage := GossipRequest{
			MessageBody: maelstrom.MessageBody{
				Type: "gossip",
			},

			Messages: b.buffer,
		}

		go func() {
			node.Send(neighbor, neighborMessage)
		}()
	}
}

func (b *Broadcaster) Stop() {
	b.stopChn <- true
}
