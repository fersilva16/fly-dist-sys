package main

import (
	"fly-dist-sys/testutils"
	"io"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestTopology1(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = testutils.NewNode()

	go main()

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(err)

	output, err := testutils.RPC(stdin, stdout, TopologyRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "topology",
			MsgID: 2,
		},

		Topology: map[string][]string{"n0": {}},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, neighbours)
}

func TestTopology2(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = testutils.NewNode()

	go main()

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0", "n1", "n2", "n3", "n4"})

	require.NoError(err)

	output, err := testutils.RPC(stdin, stdout, TopologyRequest{
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

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, neighbours)
}
