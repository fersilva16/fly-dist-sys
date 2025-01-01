package main

import (
	"encoding/json"
	test_utils "fly-dist-sys/testutils"
	"io"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

type MockClock struct{}

func (clock MockClock) Now() string {
	return "1709404427"
}

func TestGenerateSingle(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	clock = MockClock{}
	count = 0

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	})

	require.NoError(body_err)

	send_err := test_utils.Send(stdin, body)

	require.NoError(send_err)

	output, read_err := test_utils.Read(stdout)

	require.NoError(read_err)

	snaps.MatchSnapshot(t, output)
}

func TestGenerateMultiple(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	clock = MockClock{}
	count = 0

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(init_err)

	body, body_err := json.Marshal(maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	})

	require.NoError(body_err)

	send1_err := test_utils.Send(stdin, body)

	require.NoError(send1_err)

	output1, read1_err := test_utils.Read(stdout)

	require.NoError(read1_err)

	snaps.MatchSnapshot(t, output1)

	send2_err := test_utils.Send(stdin, body)

	require.NoError(send2_err)

	output2, read2_err := test_utils.Read(stdout)

	require.NoError(read2_err)

	snaps.MatchSnapshot(t, output2)
}
