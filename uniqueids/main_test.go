package main

import (
	"gossip-gloomers/testutils"

	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

type MockClock struct{}

func (clock MockClock) Now() string {
	return "1709404427"
}

func TestSingle(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()
	clock = MockClock{}
	count = 0

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	output, err := client.RPC(maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
}

func TestMultiple(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()
	clock = MockClock{}
	count = 0

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)

	go main()

	err := client.InitNode("n0", []string{"n0", "n1", "n2", "n3", "n4"})

	require.NoError(err)

	body := maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	}

	output1, err := client.RPC(body)

	require.NoError(err)

	snaps.MatchSnapshot(t, output1)

	output2, err := client.RPC(body)

	require.NoError(err)

	snaps.MatchSnapshot(t, output2)
}
