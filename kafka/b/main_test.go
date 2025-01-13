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
	kv := testutils.NewLinKV(link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	err = client.Write(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 2,
		},

		Key: "0",
		Msg: 83,
	})

	require.NoError(err)

	err = kv.HandleRead("0", float64(0))

	require.NoError(err)

	err = kv.HandleCAS("0", float64(0), float64(1), true)

	require.NoError(err)

	err = kv.HandleWrite("0-0", float64(83))

	require.NoError(err)

	sendOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, sendOutput)

	err = client.Write(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 3,
		},

		Offsets: map[string]int{
			"0": 0,
		},
	})

	require.NoError(err)

	err = kv.HandleRead("0", float64(1))

	require.NoError(err)

	err = kv.HandleRead("0-0", float64(83))

	require.NoError(err)

	pollOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, pollOutput)

	err = client.Write(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 4,
		},

		Key: "0",
		Msg: 84,
	})

	require.NoError(err)

	err = kv.HandleRead("0", float64(1))

	require.NoError(err)

	err = kv.HandleCAS("0", float64(1), float64(2), true)

	require.NoError(err)

	err = kv.HandleWrite("0-1", float64(84))

	require.NoError(err)

	sendOutput2, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, sendOutput2)

	err = client.Write(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 5,
		},

		Offsets: map[string]int{
			"0": 1,
		},
	})

	require.NoError(err)

	err = kv.HandleRead("0", float64(2))

	require.NoError(err)

	err = kv.HandleRead("0-1", float64(84))

	require.NoError(err)

	pollOutput2, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, pollOutput2)

	err = client.Write(CommitOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "commit_offsets",
			MsgID: 6,
		},

		Offsets: map[string]int{
			"0": 1,
		},
	})

	require.NoError(err)

	err = kv.HandleWrite("commit-0", float64(1))

	require.NoError(err)

	commitOffsetsOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, commitOffsetsOutput)

	err = client.Write(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 7,
		},

		Keys: []string{
			"0",
		},
	})

	require.NoError(err)

	err = kv.HandleRead("commit-0", float64(1))

	require.NoError(err)

	listCommittedOffsetOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, listCommittedOffsetOutput)
}

func TestCASConflict(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)
	kv := testutils.NewLinKV(link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	err = client.Write(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 2,
		},

		Key: "0",
		Msg: 83,
	})

	require.NoError(err)

	err = kv.HandleRead("0", float64(0))

	require.NoError(err)

	err = kv.HandleCASConflict("0", float64(0), float64(1), true)

	require.NoError(err)

	err = kv.HandleRead("0", float64(1))

	require.NoError(err)

	err = kv.HandleCAS("0", float64(1), float64(2), true)

	require.NoError(err)

	err = kv.HandleWrite("0-1", float64(83))

	require.NoError(err)

	sendOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, sendOutput)

	err = client.Write(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 3,
		},

		Offsets: map[string]int{
			"0": 0,
		},
	})

	require.NoError(err)

	err = kv.HandleRead("0", float64(2))

	require.NoError(err)

	err = kv.HandleRead("0-0", float64(82))

	require.NoError(err)

	err = kv.HandleRead("0-1", float64(83))

	require.NoError(err)

	pollOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, pollOutput)
}

func TestCommitedOffsetDoesNotExist(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)
	kv := testutils.NewLinKV(link)

	go main()

	err := client.InitNode("n0", []string{"n0"})

	require.NoError(err)

	err = client.Write(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 7,
		},

		Keys: []string{
			"0",
		},
	})

	require.NoError(err)

	err = kv.HandleReadNotExist("commit-0")

	require.NoError(err)

	listCommittedOffsetOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, listCommittedOffsetOutput)
}
