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
	kv := testutils.NewSeqKV(link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	err = client.Write(AddRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "add",
			MsgID: 2,
		},

		Delta: 1,
	})

	require.NoError(err)

	err = kv.HandleRead("counter", 0)

	require.NoError(err)

	err = kv.HandleCAS("counter", float64(0), float64(1), true)

	require.NoError(err)

	output, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, output)

	err = client.Write(maelstrom.MessageBody{
		Type:  "read",
		MsgID: 3,
	})

	require.NoError(err)

	err = kv.HandleRead("counter", 1)

	require.NoError(err)

	output, err = client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
}
