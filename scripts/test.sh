go build -o ./build ./maelstrom-echo ./maelstrom-unique-ids ./maelstrom-broadcast

maelstrom test -w echo --bin ./build/maelstrom-echo --node-count 1 --time-limit 10

maelstrom test -w unique-ids --bin ./build/maelstrom-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

maelstrom test -w broadcast --bin ./build/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10
maelstrom test -w broadcast --bin ./build/maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10
maelstrom test -w broadcast --bin ./build/maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition
