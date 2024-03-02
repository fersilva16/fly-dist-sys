module fly-dist-sys/maelstrom-echo

go 1.21.6

require github.com/jepsen-io/maelstrom/demo/go v0.0.0-20231231190402-2674df7c1076

require (
	github.com/fersilva16/fly-dist-sys/maelstrom-test-utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/gkampitakis/go-snaps v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
)

replace github.com/fersilva16/fly-dist-sys/maelstrom-test-utils => ../maelstrom-test-utils
