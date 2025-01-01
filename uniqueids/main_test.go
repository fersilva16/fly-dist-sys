package main

import (
	"fly-dist-sys/testutils"

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

	node, stdin, stdout = testutils.NewNode()
	clock = MockClock{}
	count = 0

	go main()

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(err)

	output, err := testutils.RPC(stdin, stdout, maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	})

	require.NoError(err)

	snaps.MatchSnapshot(t, output)
}

func TestGenerateMultiple(t *testing.T) {
	require := require.New(t)

	var stdin io.WriteCloser
	var stdout io.ReadCloser

	node, stdin, stdout = testutils.NewNode()
	clock = MockClock{}
	count = 0

	go main()

	err := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(err)

	body := maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	}

	output1, err := testutils.RPC(stdin, stdout, body)

	require.NoError(err)

	snaps.MatchSnapshot(t, output1)

	output2, err := testutils.RPC(stdin, stdout, body)

	require.NoError(err)

	snaps.MatchSnapshot(t, output2)
}
