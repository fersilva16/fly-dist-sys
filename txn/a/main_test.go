package main

import (
	"gossip-gloomers/testutils"
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

	txnOutput1, err := client.RPC(TxnRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "txn",
			MsgID: 2,
		},

		Txn: Txn{
			{WRITE, 1, 1},
			{READ, 2, nil},
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, txnOutput1)

	txnOutput2, err := client.RPC(TxnRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "txn",
			MsgID: 3,
		},

		Txn: Txn{
			{READ, 1, nil},
			{WRITE, 2, 2},
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, txnOutput2)
}
