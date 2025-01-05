package main

import (
	"fly-dist-sys/testutils"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
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
