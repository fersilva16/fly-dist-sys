package main

import (
	"gossip-gloomers/utils"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Broadcaster struct {
	node    *maelstrom.Node
	set     *utils.Set[int]
	stopChn chan bool
	chn     chan int
	ticker  *time.Ticker
}

func NewBroadcaster(node *maelstrom.Node, set *utils.Set[int], d time.Duration) *Broadcaster {
	return &Broadcaster{
		node:    node,
		set:     set,
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

	for {
		select {
		case <-b.chn:
		default:
			return true
		}
	}
}

func (b *Broadcaster) broadcast() {
	if !b.flush() {
		return
	}

	for _, neighbor := range node.NodeIDs() {
		if neighbor == node.ID() {
			continue
		}

		neighborMessage := GossipRequest{
			MessageBody: maelstrom.MessageBody{
				Type: "gossip",
			},

			Messages: b.set.GetAll(),
		}

		go func() {
			node.Send(neighbor, neighborMessage)
		}()
	}
}

func (b *Broadcaster) Stop() {
	b.stopChn <- true
}
