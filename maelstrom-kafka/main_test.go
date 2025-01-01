package main

import (
	"encoding/json"
	"io"
	"testing"

	test_utils "github.com/fersilva16/fly-dist-sys/maelstrom-test-utils"
	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
	"github.com/trailofbits/go-mutexasserts"
)

func TestSend1(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 0
	o.offsets = map[string]int{}
	m.messages = map[string][][]int{}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 2,
		},

		Key: "6",
		Msg: 1,
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, offsets_offset)
	snaps.MatchJSON(t, o.offsets)
	snaps.MatchJSON(t, m.messages)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestSend2(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"6": 1}
	m.messages = map[string][][]int{"6": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 2,
		},

		Key: "6",
		Msg: 1,
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, offsets_offset)
	snaps.MatchJSON(t, o.offsets)
	snaps.MatchJSON(t, m.messages)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestSend3(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"6": 1}
	m.messages = map[string][][]int{"6": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "send",
			MsgID: 2,
		},

		Key: "9",
		Msg: 1,
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
	snaps.MatchJSON(t, offsets_offset)
	snaps.MatchJSON(t, o.offsets)
	snaps.MatchJSON(t, m.messages)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestPoll1(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 0
	o.offsets = map[string]int{}
	m.messages = map[string][][]int{}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 2,
		},

		Offsets: map[string]int{},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestPoll2(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"6": 1}
	m.messages = map[string][][]int{"6": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 2,
		},

		Offsets: map[string]int{},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestPoll3(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"6": 1}
	m.messages = map[string][][]int{"6": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 2,
		},

		Offsets: map[string]int{
			"6": 1,
		},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestPoll4(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"6": 3}
	m.messages = map[string][][]int{"6": {{1, 1}, {2, 2}, {3, 3}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(PollRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "poll",
			MsgID: 2,
		},

		Offsets: map[string]int{
			"6": 2,
		},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestListCommitedOffsets1(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 0
	o.offsets = map[string]int{}
	m.messages = map[string][][]int{}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 2,
		},

		Keys: []string{"6"},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestListCommitedOffsets2(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"6": 1}
	m.messages = map[string][][]int{"6": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 2,
		},

		Keys: []string{"6"},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestListCommitedOffsets3(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 1
	o.offsets = map[string]int{"7": 1}
	m.messages = map[string][][]int{"7": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 2,
		},

		Keys: []string{"6"},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}

func TestListCommitedOffsets4(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	offsets_offset = 2
	o.offsets = map[string]int{"6": 1, "7": 1}
	m.messages = map[string][][]int{"6": {{1, 1}}, "7": {{1, 1}}}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(ListCommittedOffsetsRequest{
		MessageBody: maelstrom.MessageBody{
			Type:  "list_committed_offsets",
			MsgID: 2,
		},

		Keys: []string{"6"},
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)

	require.False(mutexasserts.RWMutexLocked(&m.mu))
	require.False(mutexasserts.RWMutexRLocked(&m.mu))
	require.False(mutexasserts.RWMutexLocked(&o.mu))
	require.False(mutexasserts.RWMutexRLocked(&o.mu))
}
