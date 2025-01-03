package testutils

import (
	"io"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Link struct {
	Node   *maelstrom.Node
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func NewLink(node *maelstrom.Node) *Link {
	inp, stdin := io.Pipe()
	stdout, outp := io.Pipe()

	node.Stdin = inp
	node.Stdout = outp

	return &Link{
		Node:   node,
		stdin:  stdin,
		stdout: stdout,
	}
}

func (l *Link) Write(src string, body any) error {
	return Write(l.stdin, src, l.Node.ID(), body)
}

func (l *Link) Read() (string, error) {
	return Read(l.stdout)
}

func (l *Link) RPC(src string, body any) (string, error) {
	return RPC(l.stdin, l.stdout, src, l.Node.ID(), body)
}
