package main

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Pair struct {
	offset, msg int
}

func (p *Pair) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int{p.offset, p.msg})
}

func (p *Pair) UnmarshalJSON(b []byte) error {
	var offsets []int

	if err := json.Unmarshal(b, &offsets); err != nil {
		return err
	}

	p.offset = offsets[0]
	p.msg = offsets[1]

	return nil
}

type NextOffsetRequest struct {
	maelstrom.MessageBody
	Key string `json:"key"`
	Msg int    `json:"msg"`
}

type NextOffsetResponse struct {
	maelstrom.MessageBody
	Offset int `json:"offset"`
}

type NewMessageRequest struct {
	maelstrom.MessageBody
	Key    string `json:"key"`
	Offset int    `json:"offset"`
	Msg    int    `json:"msg"`
}

type SendRequest struct {
	maelstrom.MessageBody
	Key string `json:"key"`
	Msg int    `json:"msg"`
}

type SendResponse struct {
	maelstrom.MessageBody
	Offset int `json:"offset"`
}

type PollRequest struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`
}

type PollResponse struct {
	maelstrom.MessageBody
	Msgs map[string][]Pair `json:"msgs"`
}

type CommitOffsetsRequest struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`
	Gossip  bool           `json:"gossip"`
}

type ListCommittedOffsetsRequest struct {
	maelstrom.MessageBody
	Keys []string `json:"keys"`
}

type ListCommittedOffsetsResponse struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`
}

var LEADER = "n0"
var node = maelstrom.NewNode()

type Messages struct {
	messages map[string][]Pair
	mu       sync.Mutex
}

func (m *Messages) Add(key string, pair Pair) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages[key] = append(m.messages[key], pair)
}

func (m *Messages) GetAfter(key string, offset int) []Pair {
	m.mu.Lock()
	defer m.mu.Unlock()

	res := []Pair{}

	for _, pair := range m.messages[key] {
		if pair.offset >= offset {
			res = append(res, pair)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].offset < res[j].offset
	})

	return res
}

type Offsets struct {
	offsets map[string]int
	mu      sync.Mutex
}

func NewOffsets() *Offsets {
	return &Offsets{
		offsets: map[string]int{},
	}
}

func (o *Offsets) Next(key string) int {
	o.mu.Lock()
	defer o.mu.Unlock()

	offset := o.offsets[key]

	o.offsets[key] += 1

	return offset
}

func (o *Offsets) Merge(m map[string]int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for key, offset := range m {
		if offset > o.offsets[key] {
			o.offsets[key] = offset
		}
	}
}

func main() {
	messages := Messages{messages: map[string][]Pair{}}
	offsets := NewOffsets()
	commitedOffsets := NewOffsets()

	node.Handle("next_offset", func(msg maelstrom.Message) error {
		var body NextOffsetRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		offset := offsets.Next(body.Key)

		messages.Add(body.Key, Pair{offset, body.Msg})

		for _, neighbor := range node.NodeIDs() {
			if neighbor == node.ID() {
				continue
			}

			if neighbor == msg.Src {
				continue
			}

			neighborMessage := NewMessageRequest{
				MessageBody: maelstrom.MessageBody{
					Type: "new_message",
				},

				Key:    body.Key,
				Offset: offset,
				Msg:    body.Msg,
			}

			go func() {
				node.Send(neighbor, neighborMessage)
			}()
		}

		resBody := NextOffsetResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "next_offset_ok",
			},

			Offset: offset,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("new_message", func(msg maelstrom.Message) error {
		var body NewMessageRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages.Add(body.Key, Pair{body.Offset, body.Msg})

		return nil
	})

	node.Handle("send", func(msg maelstrom.Message) error {
		ctx := context.Background()

		var body SendRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		var offset int

		if node.ID() == LEADER {
			offset = offsets.Next(body.Key)
		} else {
			resp, err := node.SyncRPC(ctx, LEADER, NextOffsetRequest{
				MessageBody: maelstrom.MessageBody{
					Type: "next_offset",
				},

				Key: body.Key,
				Msg: body.Msg,
			})

			if err != nil {
				return err
			}

			var respBody NextOffsetResponse

			if err := json.Unmarshal(resp.Body, &respBody); err != nil {
				return err
			}

			offset = respBody.Offset
		}

		messages.Add(body.Key, Pair{offset, body.Msg})

		if node.ID() == LEADER {
			for _, neighbor := range node.NodeIDs() {
				if neighbor == node.ID() {
					continue
				}

				neighborMessage := NewMessageRequest{
					MessageBody: maelstrom.MessageBody{
						Type: "new_message",
					},

					Key:    body.Key,
					Offset: offset,
					Msg:    body.Msg,
				}

				go func() {
					node.Send(neighbor, neighborMessage)
				}()
			}
		}

		resBody := SendResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "send_ok",
			},

			Offset: offset,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("poll", func(msg maelstrom.Message) error {
		var body PollRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgs := map[string][]Pair{}

		for key, offset := range body.Offsets {
			msgs[key] = messages.GetAfter(key, offset)
		}

		resBody := PollResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "poll_ok",
			},

			Msgs: msgs,
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body CommitOffsetsRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		commitedOffsets.Merge(body.Offsets)

		if body.Gossip {
			return nil
		}

		for _, neighbor := range node.NodeIDs() {
			if neighbor == node.ID() {
				continue
			}

			neighborMessage := CommitOffsetsRequest{
				MessageBody: maelstrom.MessageBody{
					Type: "commit_offsets",
				},

				Offsets: commitedOffsets.offsets,
				Gossip:  true,
			}

			go func() {
				node.Send(neighbor, neighborMessage)
			}()
		}

		resBody := maelstrom.MessageBody{
			Type: "commit_offsets_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body ListCommittedOffsetsRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resOffsets := map[string]int{}

		for _, key := range body.Keys {
			resOffsets[key] = commitedOffsets.offsets[key]
		}

		resBody := ListCommittedOffsetsResponse{
			MessageBody: maelstrom.MessageBody{
				Type: "list_committed_offsets_ok",
			},

			Offsets: resOffsets,
		}

		return node.Reply(msg, resBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
