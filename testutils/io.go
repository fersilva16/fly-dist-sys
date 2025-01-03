package testutils

import (
	"bufio"
	"encoding/json"
	"io"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func Write(stdin io.WriteCloser, src string, dest string, body any) error {
	bodyRaw, err := json.Marshal(body)

	if err != nil {
		return err
	}

	msg := maelstrom.Message{
		Src:  src,
		Dest: dest,
		Body: bodyRaw,
	}

	raw, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	stdin.Write(append(raw, '\n'))

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

func RPC(stdin io.WriteCloser, stdout io.ReadCloser, src string, dest string, body any) (string, error) {
	bodyRaw, err := json.Marshal(body)

	if err != nil {
		return "", err
	}

	err = Write(stdin, src, dest, bodyRaw)

	if err != nil {
		return "", err
	}

	msg, err := Read(stdout)

	if err != nil {
		return "", err
	}

	return msg, nil
}
