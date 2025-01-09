package main

import "sync"

type CRDT struct {
	m  map[string]int
	mu sync.RWMutex
}

func NewCRDT() *CRDT {
	return &CRDT{
		m: map[string]int{},
	}
}

func (c *CRDT) Sync(src string, count int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[src] = count
}

func (c *CRDT) Read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0

	for _, value := range c.m {
		count += value
	}

	return count
}
