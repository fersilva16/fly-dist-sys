package test_utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func Send(stdin io.WriteCloser, body json.RawMessage) error {
	msg := maelstrom.Message{
		Src: "c0",
		Dest: "n0",
		Body: body,
	};

	msgJSON, err := json.Marshal(msg);

	if err != nil {
		return err;
	}

	stdin.Write(msgJSON);
	stdin.Write([]byte{'\n'});
	
	return nil;
}

func Read(stdout io.ReadCloser) (string, error) {
	in := bufio.NewScanner(stdout);

	for in.Scan() {
		return in.Text(), nil;
	}

	if err := in.Err(); err != nil {
		return "", err;
	}

	return "", nil;
}

func InitNode(stdin io.WriteCloser, stdout io.ReadCloser, node_id string) error {
	body, body_err := json.Marshal(maelstrom.InitMessageBody{
		MessageBody: maelstrom.MessageBody{
			Type: "init",
			MsgID: 1,
		},

		NodeID: node_id,
		NodeIDs: []string{ node_id },
	});

	if body_err != nil {
		return body_err;
	}

	err := Send(stdin, body);

	if err != nil {
		return err;
	}

	msg, read_err := Read(stdout);

	if read_err != nil {
		return read_err;
	}

	var message maelstrom.Message;
	
	if err := json.Unmarshal([]byte(msg), &message); err != nil {
		return err;
	};
	
	var msg_body maelstrom.MessageBody;
	
	if err := json.Unmarshal(message.Body, &msg_body); err != nil {
		return err;
	};

	if msg_body.Type != "init_ok" || msg_body.InReplyTo != 1 {
		return fmt.Errorf("invalid message: %s", message.Body);
	}
	
	return nil;
}

func NewNode() (*maelstrom.Node, io.WriteCloser, io.ReadCloser) {
	node := maelstrom.NewNode();

	inp, stdin := io.Pipe();
	stdout, outp := io.Pipe();
	
	node.Stdin = inp;
	node.Stdout = outp;

	return node, stdin, stdout;
}
