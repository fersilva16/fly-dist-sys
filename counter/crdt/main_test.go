package main

import (
	"gossip-gloomers/testutils"
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
	node1 := testutils.NewClient("n1", link)

	go main()

	err := client.InitNode("n0", []string{"n0", "n1"})

	require.NoError(err)

	err = client.Write(AddRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "add",
			MsgID: 2,
		},

		Delta: 1,
	})

	require.NoError(err)

	output, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, output)

	propagateOutput, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, propagateOutput)

	err = node1.Write(PropagateRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "propagate",
			MsgID: 4,
		},

		Count: 1,
	})

	require.NoError(err)

	err = client.Write(maelstrom.MessageBody{
		Type:  "read",
		MsgID: 3,
	})

	require.NoError(err)

	output, err = client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
}
