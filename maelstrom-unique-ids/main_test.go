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

type MockClock struct{}

func (clock MockClock) Now() string {
	return "1709404427"
}

func TestGenerate1(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	clock = MockClock{}

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

func TestGenerate2(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	clock = MockClock{}
	count = 1

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

func TestGenerate3(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = test_utils.NewNode()
	clock = MockClock{}

	go main()

	init_err := test_utils.InitNode(stdin, stdout, "n1", []string{"n1"})

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
