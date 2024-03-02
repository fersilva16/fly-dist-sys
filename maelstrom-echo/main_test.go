package main

import (
	"encoding/json"
	"os/exec"
	"testing"

	test_utils "github.com/fersilva16/fly-dist-sys/maelstrom-test-utils"
	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	require := require.New(t);

	cmd := exec.Command("go", "run", "./main.go");

	stdin, _ := cmd.StdinPipe();
	stdout, _ := cmd.StdoutPipe();

	if err := cmd.Start(); err != nil {
		t.Error(err);

		return;
	}

	err := test_utils.InitNode(stdin, stdout);

	require.NoError(err);

	body, body_err := json.Marshal(EchoRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "echo",
			MsgID: 2,
		},

		Echo: "Please echo 1",
	});

	require.NoError(body_err);

	send_err := test_utils.Send(stdin,body);

	require.NoError(send_err);

	output, read_err := test_utils.Read(stdout);

	require.NoError(read_err);

	snaps.MatchSnapshot(t, output);

	if err := cmd.Process.Kill(); err != nil {
		t.Error(err)
	}
}
