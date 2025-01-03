package testutils

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Client struct {
	id   string
	link *Link
}

func NewClient(id string, link *Link) *Client {
	return &Client{
		id:   id,
		link: link,
	}
}

func (c *Client) Write(body any) error {
	return c.link.Write(c.id, body)
}

func (c *Client) Read() (string, error) {
	output, err := c.link.Read()

	if err != nil {
		return "", err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return "", err
	}

	if msg.Dest != c.id {
		return "", fmt.Errorf("invalid dest for client %s: %s", c.id, msg.Dest)
	}

	return output, nil
}

func (c *Client) RPC(body any) (string, error) {
	err := c.Write(body)

	if err != nil {
		return "", err
	}

	output, err := c.Read()

	if err != nil {
		return "", err
	}

	return output, nil
}

func (c *Client) InitNode(nodeID string, nodeIDs []string) error {
	output, err := c.RPC(maelstrom.InitMessageBody{
		MessageBody: maelstrom.MessageBody{
			Type:  "init",
			MsgID: 1,
		},

		NodeID:  nodeID,
		NodeIDs: nodeIDs,
	})

	if err != nil {
		return err
	}

	var msg maelstrom.Message

	if err := json.Unmarshal([]byte(output), &msg); err != nil {
		return err
	}

	var body maelstrom.MessageBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Type != "init_ok" || body.InReplyTo != 1 {
		return fmt.Errorf("invalid message: %s", msg.Body)
	}

	return nil
}
