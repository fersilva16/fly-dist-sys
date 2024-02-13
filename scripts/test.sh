go build -o build
maelstrom test -w echo --bin ./build/fly-dist-sys --node-count 1 --time-limit 10