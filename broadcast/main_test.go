package main

import (
	"fly-dist-sys/testutils"
	"io"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestSingle(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = testutils.NewNode()

	go main()

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(err)

	topologyOutput, err := testutils.RPC(stdin, stdout, TopologyRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "topology",
			MsgID: 2,
		},

		Topology: map[string][]string{"n0": {}},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, topologyOutput)

	broadcastOutput, err := testutils.RPC(stdin, stdout, BroadcastRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "broadcast",
			MsgID: 2,
		},

		Message: 1,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutput)

	readOutput, err := testutils.RPC(stdin, stdout, maelstrom.MessageBody{
		Type:  "read",
		MsgID: 2,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, readOutput)
}
