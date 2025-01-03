package main

import (
	"fly-dist-sys/testutils"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestSingle(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	topologyOutput, err := client.RPC(TopologyRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "topology",
			MsgID: 2,
		},

		Topology: map[string][]string{"n0": {}},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, topologyOutput)

	broadcastOutput, err := client.RPC(BroadcastRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "broadcast",
			MsgID: 2,
		},

		Message: 1,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutput)

	readOutput, err := client.RPC(maelstrom.MessageBody{
		Type:  "read",
		MsgID: 2,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, readOutput)
}

func TestMulti(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)

	go main()

	err := client.InitNode("n0", []string{"n0", "n1", "n2", "n3", "n4"})

	require.NoError(err)

	topologyOutput, err := client.RPC(TopologyRequest{
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

	broadcastOutput, err := client.RPC(BroadcastRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "broadcast",
			MsgID: 3,
		},

		Message: 1,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutput)

	broadcastOutputN3, err := link.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutputN3)

	broadcastOutputN1, err := link.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, broadcastOutputN1)

	readOutput, err := client.RPC(maelstrom.MessageBody{
		Type:  "read",
		MsgID: 4,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, readOutput)
}
