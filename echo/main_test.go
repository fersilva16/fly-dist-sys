package main

import (
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

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(err)

	output, err := testutils.RPC(stdin, stdout, EchoRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "echo",
			MsgID: 2,
		},

		Echo: "Please echo 1",
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
}
