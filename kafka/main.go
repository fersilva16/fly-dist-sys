package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type SendRequest struct {
	maelstrom.MessageBody
	Key string `json:"key"`
	Msg int    `json:"msg"`
}

type PollRequest struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`
}

type CommitOffsetsRequest struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`
}

type ListCommittedOffsetsRequest struct {
	maelstrom.MessageBody
	Keys []string `json:"keys"`
}

type Offsets struct {
	mu      sync.RWMutex
	offsets map[string]int
}

type Messages struct {
	mu       sync.RWMutex
	messages map[string][][]int
}

var node = maelstrom.NewNode()
var offsetsOffset = 0
var o = Offsets{offsets: make(map[string]int)}
var m = Messages{messages: make(map[string][][]int)}

func main() {
	node.Handle("send", func(msg maelstrom.Message) error {
		var body SendRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		o.mu.Lock()

		if o.offsets[body.Key] == 0 {
			o.offsets[body.Key] = offsetsOffset * 1000
			offsetsOffset++
		}

		o.offsets[body.Key]++

		offset := o.offsets[body.Key]

		o.mu.Unlock()

		m.mu.Lock()

		if m.messages[body.Key] == nil {
			m.messages[body.Key] = [][]int{}
		}

		m.messages[body.Key] = append(m.messages[body.Key], []int{offset, body.Msg})

		m.mu.Unlock()

		resBody := map[string]any{
			"type":   "send_ok",
			"offset": offset,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("poll", func(msg maelstrom.Message) error {
		var body PollRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resMessages := make(map[string][][]int)

		for key, offset := range body.Offsets {
			m.mu.RLock()

			msgs := m.messages[key]

			m.mu.RUnlock()

			if msgs == nil {
				continue
			}

			resMessages[key] = [][]int{}

			for i := 0; i < len(msgs); i++ {
				if msgs[i][0] < offset {
					continue
				}

				resMessages[key] = append(resMessages[key], msgs[i])
			}
		}

		resBody := map[string]any{
			"type": "poll_ok",
			"msgs": resMessages,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body CommitOffsetsRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		for key, offset := range body.Offsets {
			o.mu.Lock()

			o.offsets[key] = offset

			o.mu.Unlock()
		}

		resBody := map[string]any{
			"type": "commit_offsets_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body ListCommittedOffsetsRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resOffsets := map[string]int{}

		for i := 0; i < len(body.Keys); i++ {
			o.mu.RLock()

			key := body.Keys[i]
			offset := o.offsets[key]

			o.mu.RUnlock()

			if offset == 0 {
				continue
			}

			resOffsets[key] = offset
		}

		resBody := map[string]any{
			"type":    "list_committed_offsets_ok",
			"offsets": resOffsets,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
