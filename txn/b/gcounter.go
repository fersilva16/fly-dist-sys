package main

import (
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type GCounter struct {
	node *maelstrom.Node
	m    map[string]int
	mu   sync.Mutex
}

func NewGCounter(node *maelstrom.Node) *GCounter {
	return &GCounter{
		node: node,
		m:    map[string]int{},
	}
}

func (c *GCounter) Sync(src string, count int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[src] = count
}

func (c *GCounter) Increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[c.node.ID()]++

	return c.read()
}

func (c *GCounter) GetLocal() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.m[c.node.ID()]
}

func (c *GCounter) read() int {
	count := 0

	for _, value := range c.m {
		count += value
	}

	return count
}
