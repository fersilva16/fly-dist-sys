go build -o build
# maelstrom test -w echo --bin ./build/fly-dist-sys --node-count 1 --time-limit 10
# maelstrom test -w unique-ids --bin ./build/fly-dist-sys --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
# maelstrom test -w broadcast --bin ./build/fly-dist-sys --node-count 1 --time-limit 20 --rate 10
maelstrom test -w broadcast --bin ./build/fly-dist-sys --node-count 5 --time-limit 20 --rate 10 --log-stderr