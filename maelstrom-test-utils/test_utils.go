package test_utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

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

func InitNode(stdin io.WriteCloser, stdout io.ReadCloser) error {
	body, body_err := json.Marshal(maelstrom.InitMessageBody{
		MessageBody: maelstrom.MessageBody{
			Type: "init",
			MsgID: 1,
		},

		NodeID: "n0",
		NodeIDs: []string{"n0"},
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

func NewNode() (*exec.Cmd, io.WriteCloser, io.ReadCloser, error) {
	cmd := exec.Command("go", "run", "./main.go");

	stdin, _ := cmd.StdinPipe();
	stdout, _ := cmd.StdoutPipe();

	if err := cmd.Start(); err != nil {
		return cmd, stdin, stdout, err;
	}

	if err := InitNode(stdin, stdout); err != nil {
		return cmd, stdin, stdout, err;
	}

	return cmd, stdin, stdout, nil;
}