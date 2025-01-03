package main

import (
	"fly-dist-sys/testutils"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	output, err := client.RPC(EchoRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "echo",
			MsgID: 2,
		},

		Echo: "Please echo 1",
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
}
