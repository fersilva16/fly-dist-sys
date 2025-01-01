package main

import (
	"encoding/json"
	"fly-dist-sys/testutils"
	"io"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = testutils.NewNode()

	go main()

	init_err := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(EchoRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "echo",
			MsgID: 2,
		},

		Echo: "Please echo 1",
	})

	require.NoError(body_err)

	send_err := testutils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := testutils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
}
