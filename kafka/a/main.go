package main

import (
	"encoding/json"
	"log"
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
}

type ListCommittedOffsetsRequest struct {
	maelstrom.MessageBody
	Keys []string `json:"keys"`
}

type ListCommittedOffsetsResponse struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`
}

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

	return res
}

type Offsets struct {
	offsets map[string]int
	mu      sync.Mutex
}

func (o *Offsets) Get(key string) int {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.offsets[key]
}

func (o *Offsets) Merge(m map[string]int) int {
	o.mu.Lock()
	defer o.mu.Unlock()

	return 0
}

func main() {
	messages := Messages{messages: map[string][]Pair{}}
	offsets := map[string]int{}
	commitedOffsets := map[string]int{}

	node.Handle("send", func(msg maelstrom.Message) error {
		var body SendRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		offset := offsets[body.Key]

		messages.Add(body.Key, Pair{offset, body.Msg})

		offsets[body.Key] += 1

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

		for key, offset := range body.Offsets {
			commitedOffsets[key] = offset
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
			resOffsets[key] = commitedOffsets[key]
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
