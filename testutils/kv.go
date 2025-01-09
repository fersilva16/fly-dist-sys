package testutils

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type KVReadMessageBody struct {
	maelstrom.MessageBody
	Key string `json:"key"`
}

type KVReadOKMessageBody struct {
	maelstrom.MessageBody
	Value any `json:"value"`
}

type KVWriteMessageBody struct {
	maelstrom.MessageBody
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type KVCASMessageBody struct {
	maelstrom.MessageBody
	Key               string `json:"key"`
	From              any    `json:"from"`
	To                any    `json:"to"`
	CreateIfNotExists bool   `json:"create_if_not_exists,omitempty"`
}

type KV struct {
	typ  string
	link *Link
}

func NewKV(typ string, link *Link) *KV {
	return &KV{
		typ:  typ,
		link: link,
	}
}

func NewLinKV(link *Link) *KV {
	return NewKV(maelstrom.LinKV, link)
}

func NewSeqKV(link *Link) *KV {
	return NewKV(maelstrom.SeqKV, link)
}

func NewLWWKV(link *Link) *KV {
	return NewKV(maelstrom.LWWKV, link)
}

func (kv *KV) Write(body any) error {
	return kv.link.Write(kv.typ, body)
}

func (kv *KV) Read() (string, error) {
	output, err := kv.link.Read()

	if err != nil {
		return "", err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return "", err
	}

	if msg.Dest != kv.typ {
		return "", fmt.Errorf("invalid dest for KV %s: %s", kv.typ, msg.Dest)
	}

	return output, nil
}

func (kv *KV) HandleRead(key string, value any) error {
	output, err := kv.Read()

	if err != nil {
		return err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return err
	}

	var body KVReadMessageBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Type != "read" {
		return fmt.Errorf("invalid type: %s (expected read)", body.Type)
	}

	if body.Key != key {
		return fmt.Errorf("invalid key: %s (expected %s)", body.Key, key)
	}

	err = kv.Write(KVReadOKMessageBody{
		MessageBody: maelstrom.MessageBody{
			Type:      "read_ok",
			InReplyTo: body.MsgID,
		},

		Value: value,
	})

	if err != nil {
		return err
	}

	return nil
}

func (kv *KV) HandleWrite(key string, value any) error {
	output, err := kv.Read()

	if err != nil {
		return err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return err
	}

	var body KVWriteMessageBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Type != "write" {
		return fmt.Errorf("invalid type: %s (expected write)", body.Type)
	}

	if body.Key != key {
		return fmt.Errorf("invalid key: %s (expected %s)", body.Key, key)
	}

	if body.Value != value {
		return fmt.Errorf("invalid value: %v (expected %v)", body.Value, value)
	}

	err = kv.Write(maelstrom.MessageBody{
		Type:      "write_ok",
		InReplyTo: body.MsgID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (kv *KV) HandleCAS(key string, from, to any, createIfNotExists bool) error {
	output, err := kv.Read()

	if err != nil {
		return err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return err
	}

	var body KVCASMessageBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Type != "cas" {
		return fmt.Errorf("invalid type: %s (expected cas)", body.Type)
	}

	if body.Key != key {
		return fmt.Errorf("invalid key: %s (expected %s)", body.Key, key)
	}

	if body.From != from {
		return fmt.Errorf("invalid from: %v (expected %v)", body.From, from)
	}

	if body.To != to {
		return fmt.Errorf("invalid to: %v (expected %v)", body.To, to)
	}

	if body.CreateIfNotExists != createIfNotExists {
		return fmt.Errorf("invalid create_if_not_exists: %t (expected %t)", body.CreateIfNotExists, createIfNotExists)
	}

	err = kv.Write(maelstrom.MessageBody{
		Type:      "cas_ok",
		InReplyTo: body.MsgID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (kv *KV) HandleCASConflict(key string, from, to any, createIfNotExists bool) error {
	output, err := kv.Read()

	if err != nil {
		return err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return err
	}

	var body KVCASMessageBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Type != "cas" {
		return fmt.Errorf("invalid type: %s (expected cas)", body.Type)
	}

	if body.Key != key {
		return fmt.Errorf("invalid key: %s (expected %s)", body.Key, key)
	}

	if body.From != from {
		return fmt.Errorf("invalid from: %v (expected %v)", body.From, from)
	}

	if body.To != to {
		return fmt.Errorf("invalid to: %v (expected %v)", body.To, to)
	}

	if body.CreateIfNotExists != createIfNotExists {
		return fmt.Errorf("invalid create_if_not_exists: %t (expected %t)", body.CreateIfNotExists, createIfNotExists)
	}

	err = kv.Write(maelstrom.MessageBody{
		Type:      "error",
		Code:      22,
		Text:      fmt.Sprintf("current value MOCKED is not %d", body.From),
		InReplyTo: body.MsgID,
	})

	if err != nil {
		return err
	}

	return nil
}
