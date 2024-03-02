package main

import (
	"bufio"
	"encoding/json"
	"io"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func send(stdin io.WriteCloser, body json.RawMessage) error {
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

func read(stdout io.ReadCloser) (string, error) {
	in := bufio.NewScanner(stdout);

	for in.Scan() {
		return in.Text(), nil;
	}

	if err := in.Err(); err != nil {
		return "", err;
	}

	return "", nil;
}

func init_node(stdin io.WriteCloser, stdout io.ReadCloser) error {
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

	err := send(stdin, body);

	if err != nil {
		return err;
	}


	in := bufio.NewScanner(stdout);

	for in.Scan() {
		var message maelstrom.Message;
		
		if err := json.Unmarshal(in.Bytes(), &message); err != nil {
			return err;
		};
		
		var body maelstrom.MessageBody;
		
		if err := json.Unmarshal(message.Body, &body); err != nil {
			return err;
		};

		if body.Type == "init_ok" && body.InReplyTo == 1 {
			return nil;
		}
	}

	if err := in.Err(); err != nil {
		return err;
	}
	
	return nil;
}