package main

import (
	"encoding/json"
	test_utils "fly-dist-sys/testutils"
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

	node, stdin, stdout = test_utils.NewNode()

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(TopologyRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "topology",
			MsgID: 2,
		},

		Topology: map[string][]string{"n0": {}},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, neighbours)
}

func TestTopology2(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0", "n1", "n2", "n3", "n4"})

	require.NoError(init_err)

	body, body_err := json.Marshal(TopologyRequest{
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

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, neighbours)
}
