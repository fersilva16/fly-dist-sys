package main

import (
	"encoding/json"
	"io"
	"testing"

	test_utils "github.com/fersilva16/fly-dist-sys/maelstrom-test-utils"
	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestSend1(t *testing.T) {
	require := require.New(t);

	var stdin io.WriteCloser;
	var stdout io.ReadCloser;

	node, stdin, stdout = test_utils.NewNode();

	go main();
	
	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{ "n0" });

	require.NoError(init_err);

	body, body_err := json.Marshal(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "send",
			MsgID: 2,
		},

		Key: "6",
		Msg: 1,
	});

	require.NoError(body_err);

	send_err := test_utils.Send(stdin,body);

	require.NoError(send_err);

	output, read_err := test_utils.Read(stdout);

	require.NoError(read_err);

	snaps.MatchSnapshot(t, output);
	snaps.MatchJSON(t, offsets_offset)
	snaps.MatchJSON(t, o.offsets)
	snaps.MatchJSON(t, m.messages)
}

func TestSend2(t *testing.T) {
	require := require.New(t);

	var stdin io.WriteCloser;
	var stdout io.ReadCloser;

	node, stdin, stdout = test_utils.NewNode();
	offsets_offset = 1;
	o.offsets = map[string]int{ "6": 1 };
	m.messages = map[string][][]int{ "6": { { 1, 1 } } };

	go main();
	
	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{ "n0" });

	require.NoError(init_err);

	body, body_err := json.Marshal(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "send",
			MsgID: 2,
		},

		Key: "6",
		Msg: 1,
	});

	require.NoError(body_err);

	send_err := test_utils.Send(stdin,body);

	require.NoError(send_err);

	output, read_err := test_utils.Read(stdout);

	require.NoError(read_err);

	snaps.MatchSnapshot(t, output);
	snaps.MatchJSON(t, offsets_offset)
	snaps.MatchJSON(t, o.offsets)
	snaps.MatchJSON(t, m.messages)
}

func TestSend3(t *testing.T) {
	require := require.New(t);

	var stdin io.WriteCloser;
	var stdout io.ReadCloser;

	node, stdin, stdout = test_utils.NewNode();
	offsets_offset = 1;
	o.offsets = map[string]int{ "6": 1 };
	m.messages = map[string][][]int{ "6": { { 1, 1 } } };

	go main();
	
	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{ "n0" });

	require.NoError(init_err);

	body, body_err := json.Marshal(SendRequest{
		MessageBody: maelstrom.MessageBody{
			Type: "send",
			MsgID: 2,
		},

		Key: "9",
		Msg: 1,
	});

	require.NoError(body_err);

	send_err := test_utils.Send(stdin,body);

	require.NoError(send_err);

	output, read_err := test_utils.Read(stdout);

	require.NoError(read_err);

	snaps.MatchSnapshot(t, output);
	snaps.MatchJSON(t, offsets_offset)
	snaps.MatchJSON(t, o.offsets)
	snaps.MatchJSON(t, m.messages)
}

