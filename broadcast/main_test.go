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

func TestMulti(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = testutils.NewNode()

	go main()

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0", "n1", "n2", "n3", "n4"})

	require.NoError(err)

	topologyOutput, err := testutils.RPC(stdin, stdout, TopologyRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "topology",
			MsgID: 2,
		},

		Topology: map[string][]string{
			"n0": {"n3", "n1"},
			"n1": {"n4", "n2", "n0"},
			"n2": {"n1"},
			"n3": {"n0", "n4"},
			"n4": {"n1", "n3"},
		},
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

	broadcastOutputN3, err := testutils.Read(stdout)

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutputN3)

	broadcastOutputN1, err := testutils.Read(stdout)

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutputN1)

	readOutput, err := testutils.RPC(stdin, stdout, maelstrom.MessageBody{
		Type:  "read",
		MsgID: 2,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, readOutput)
}
