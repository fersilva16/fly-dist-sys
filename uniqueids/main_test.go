package main

import (
	"encoding/json"
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

	initErr := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(initErr)

	body, bodyErr := json.Marshal(maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	})

	require.NoError(bodyErr)

	sendErr := testutils.Send(stdin, body)

	require.NoError(sendErr)

	output, readErr := testutils.Read(stdout)

	require.NoError(readErr)

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

	initErr := testutils.InitNode(stdin, stdout, "n0", []string{"n0"})

	require.NoError(initErr)

	body, bodyErr := json.Marshal(maelstrom.MessageBody{
		Type:  "generate",
		MsgID: 2,
	})

	require.NoError(bodyErr)

	send1Err := testutils.Send(stdin, body)

	require.NoError(send1Err)

	output1, read1Err := testutils.Read(stdout)

	require.NoError(read1Err)

	snaps.MatchSnapshot(t, output1)

	send2Err := testutils.Send(stdin, body)

	require.NoError(send2Err)

	output2, read2Err := testutils.Read(stdout)

	require.NoError(read2Err)

	snaps.MatchSnapshot(t, output2)
}
