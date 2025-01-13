package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

func main() {
	kv := maelstrom.NewLinKV(node)

	node.Handle("send", func(msg maelstrom.Message) error {
		ctx := context.Background()

		var body SendRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		var offset int

		for {
			err := kv.ReadInto(ctx, body.Key, &offset)

			if err != nil {
				offset = 0
			}

			err = kv.CompareAndSwap(ctx, body.Key, offset, offset+1, true)

			if err == nil {
				break
			}
		}

		msgKvKey := fmt.Sprintf("%s-%d", body.Key, offset)
		err := kv.Write(ctx, msgKvKey, body.Msg)

		if err != nil {
			return err
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
		ctx := context.Background()

		var body PollRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgs := map[string][]Pair{}

		for key, offset := range body.Offsets {
			currentOffset, err := kv.ReadInt(ctx, key)

			if err != nil {
				continue
			}

			for i := offset; i <= currentOffset; i++ {
				msg, err := kv.ReadInt(ctx, fmt.Sprintf("%s-%d", key, i))

				if err != nil {
					continue
				}

				msgs[key] = append(msgs[key], Pair{i, msg})
			}
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
		ctx := context.Background()

		var body CommitOffsetsRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		for key, offset := range body.Offsets {
			kvKey := fmt.Sprintf("commit-%s", key)
			err := kv.Write(ctx, kvKey, offset)

			if err != nil {
				return err
			}
		}

		resBody := maelstrom.MessageBody{
			Type: "commit_offsets_ok",
		}

		return node.Reply(msg, resBody)
	})

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		ctx := context.Background()

		var body ListCommittedOffsetsRequest

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resOffsets := map[string]int{}

		for _, key := range body.Keys {
			kvKey := fmt.Sprintf("commit-%s", key)
			offset, err := kv.ReadInt(ctx, kvKey)

			if err != nil {
				offset = 0
			}

			resOffsets[key] = offset
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
