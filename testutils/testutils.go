package testutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func Send(stdin io.WriteCloser, body json.RawMessage) error {
	msg := maelstrom.Message{
		Src:  "c0",
		Dest: "n0",
		Body: body,
	}

	msgJSON, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	stdin.Write(append(msgJSON, '\n'))

	return nil
}

func Read(stdout io.ReadCloser) (string, error) {
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		return in.Text(), nil
	}

	if err := in.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func InitNode(stdin io.WriteCloser, stdout io.ReadCloser, nodeId string, nodeIds []string) error {
	body, bodyErr := json.Marshal(maelstrom.InitMessageBody{
		MessageBody: maelstrom.MessageBody{
			Type:  "init",
			MsgID: 1,
		},

		NodeID:  nodeId,
		NodeIDs: nodeIds,
	})

	if bodyErr != nil {
		return bodyErr
	}

	err := Send(stdin, body)

	if err != nil {
		return err
	}

	msg, readErr := Read(stdout)

	if readErr != nil {
		return readErr
	}

	var message maelstrom.Message

	if err := json.Unmarshal([]byte(msg), &message); err != nil {
		return err
	}

	var msgBody maelstrom.MessageBody

	if err := json.Unmarshal(message.Body, &msgBody); err != nil {
		return err
	}

	if msgBody.Type != "init_ok" || msgBody.InReplyTo != 1 {
		return fmt.Errorf("invalid message: %s", message.Body)
	}

	return nil
}

func NewNode() (*maelstrom.Node, io.WriteCloser, io.ReadCloser) {
	node := maelstrom.NewNode()

	inp, stdin := io.Pipe()
	stdout, outp := io.Pipe()

	node.Stdin = inp
	node.Stdout = outp

	return node, stdin, stdout
}
