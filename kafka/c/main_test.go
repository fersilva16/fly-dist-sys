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

	sendOutput, err := client.RPC(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 2,
		},

		Key: "0",
		Msg: 83,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, sendOutput)

	newMessageOutput, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, newMessageOutput)

	pollOutput, err := client.RPC(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 3,
		},

		Offsets: map[string]int{
			"0": 0,
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, pollOutput)

	sendOutput2, err := client.RPC(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 4,
		},

		Key: "0",
		Msg: 84,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, sendOutput2)

	newMessageOutput2, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, newMessageOutput2)

	pollOutput2, err := client.RPC(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 5,
		},

		Offsets: map[string]int{
			"0": 1,
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, pollOutput2)

	commitOffsetsOutput, err := client.RPC(CommitOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "commit_offsets",
			MsgID: 6,
		},

		Offsets: map[string]int{
			"0": 1,
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, commitOffsetsOutput)

	commitOffsetsOutput1, err := node1.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, commitOffsetsOutput1)

	listCommittedOffsetOutput, err := client.RPC(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 7,
		},

		Keys: []string{
			"0",
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, listCommittedOffsetOutput)
}

func TestFollower(t *testing.T) {
	require := require.New(t)

	node = maelstrom.NewNode()

	link := testutils.NewLink(node)
	client := testutils.NewClient("c0", link)
	node0 := testutils.NewClient("n0", link)

	go main()

	err := client.InitNode("n1", []string{"n0", "n1"})

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

	nextOffsetOutput, err := node0.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, nextOffsetOutput)

	err = node0.Write(NextOffsetResponse{
		MessageBody: maelstrom.MessageBody{
			Type:      "next_offset_ok",
			InReplyTo: 1,
		},

		Offset: 0,
	})

	require.NoError(err)

	sendOutput, err := client.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, sendOutput)

	err = node0.Write(NewMessageRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "new_message",
		},

		Key:    "0",
		Offset: 1,
		Msg:    84,
	})

	require.NoError(err)

	pollOutput, err := client.RPC(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 4,
		},

		Offsets: map[string]int{
			"0": 0,
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, pollOutput)

	commitOffsetsOutput, err := client.RPC(CommitOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "commit_offsets",
			MsgID: 6,
		},

		Offsets: map[string]int{
			"0": 1,
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, commitOffsetsOutput)

	commitOffsetsOutput1, err := node0.Read()

	require.NoError(err)

	snaps.MatchSnapshot(t, commitOffsetsOutput1)

	listCommittedOffsetOutput, err := client.RPC(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 7,
		},

		Keys: []string{
			"0",
		},
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, listCommittedOffsetOutput)
}
