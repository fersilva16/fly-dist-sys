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

	initErr := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(initErr)

	body, bodyErr := json.Marshal(EchoRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "echo",
			MsgID: 2,
		},

		Echo: "Please echo 1",
	})

	require.NoError(bodyErr)

	sendErr := testutils.Send(stdin, body)

	require.NoError(sendErr)

	output, readErr := testutils.Read(stdout)

	require.NoError(readErr)

	snaps.MatchSnapshot(t, output)
}
