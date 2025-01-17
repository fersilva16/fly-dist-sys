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

	replicateOutput1, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, replicateOutput1)

	err = node1.Write(ReplicateResponse{
		MessageBody: maelstrom.MessageBody{
			Type: "replicate_ok",
		},

		Keys: []int{1},
	})

	require.NoError(err)

	err = node1.Write(ReplicateRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "replicate",
		},

		Clock: 1,
		Snapshot: map[int]Value{
			2: {
				Value: 2,
				TxnId: 1,
			},
		},
	})

	require.NoError(err)

	replicateOutput2, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, replicateOutput2)

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

	replicateOutput3, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, replicateOutput3)

	err = node1.Write(ReplicateResponse{
		MessageBody: maelstrom.MessageBody{
			Type: "replicate_ok",
		},

		Keys: []int{2},
	})

	require.NoError(err)
}

func TestPartition(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)
	node1 := testutils.NewClient("n1", link)

	go main()

	err := client.InitNode("n0", []string{"n0", "n1"})

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

	replicateOutput1, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, replicateOutput1)

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

	replicateOutput2, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, replicateOutput2)

	err = node1.Write(ReplicateResponse{
		MessageBody: maelstrom.MessageBody{
			Type: "replicate_ok",
		},

		Keys: []int{1, 2},
	})

	require.NoError(err)
}
