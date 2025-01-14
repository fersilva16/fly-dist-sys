# Gossip Glomers

Solutions for the Distributed System Challenges from [Fly.io](https://fly.io) and [Kyle Kingsbury](https://aphyr.com/about)

## Solutions

### [Challenge #1: Echo](https://fly.io/dist-sys/1)

Hello world! Just echoes back what you send it.

[solution](echo/main.go) / [tests](echo/main_test.go)

### [Challenge #2: Unique IDs](https://fly.io/dist-sys/2)

Generates a unique ID for each node. Inspired by [MongoDB's ObjectId](https://www.mongodb.com/docs/manual/reference/method/ObjectId).

Each ID is a concatenation of:

- Node ID: to avoid collisions with other nodes
- Current time: to avoid collisions with messages sent at the same time
- Count: to avoid collisions with messages sent at the same time and node

[solution](uniqueids/main.go) / [tests](uniqueids/main_test.go)

### [Challenge #3a: Single-Node Broadcast](https://fly.io/dist-sys/3a)

A basic node that receives messages, saves them to a slice, and returns on the `read` message.

[solution](broadcast/a/main.go) / [tests](broadcast/a/main_test.go)

### [Challenge #3b: Multi-Node Broadcast](https://fly.io/dist-sys/3b)

Same as the previous challenge, but now it broadcast the message to all neighbors.

[solution](broadcast/b/main.go) / [tests](broadcast/b/main_test.go)

### [Challenge #3c: Fault Tolerant Broadcast](https://fly.io/dist-sys/3c)

Node spawns a new thread for each message sent to a neighbor, and exponentially backoff if the neighbor doesn't respond.

This solution is pretty dumb and not the most efficient one, here's other ideas of how it could be done:

- Use a map to keep track of which messages are not yet acknowledged and retry them on the next messages
- Use a separated thread with a ticker to periodically distribute the messages to neighbors

[solution](broadcast/c/main.go) / [tests](broadcast/c/main_test.go)

### [Challenge #3d: Efficient Broadcast, Part I](https://fly.io/dist-sys/3d)

The node receiving the message broadcasts to all the other nodes in network with a send-and-forget approach with the `gossip` message type.

Results from a run:

```
All checks passed
  Messages per op: 11.789474/30
  Median latency: 84/400
  Maximum latency: 105/600
```

[solution](broadcast/d/main.go) / [tests](broadcast/d/main_test.go)

### [Challenge #3e: Efficient Broadcast, Part II](https://fly.io/dist-sys/3e)

The previous solution also works for this challenge, but I've taken a step further and tried to make it send as less messages per operation as possible.

It's the same idea as #3d, but it spawns a new thread that broadcasts the messages every 1.5s with a buffered channel and only broadcasts if there are new messages.

Results:

```
All checks passed
  Messages per op: 4.728111/20
  Median latency: 791/1000
  Maximum latency: 1584/2000
```

[solution](broadcast/e/main.go) / [tests](broadcast/e/main_test.go)

### [Challenge #4: Grow-Only Counter](https://fly.io/dist-sys/4)

I've made 2 solutions for this challenge, one that uses a OT approach and another with a CRDT solution.

OT (Operational Transformation) solution uses the Sequential KV store from Maelstrom and tries to Compare and Swap the value of the counter.

[solution](counter/ot/main.go) / [tests](counter/ot/main_test.go)

CRDT (Conflict-free Replicated Data Type) solution uses a map to store the counter of the other nodes and propagates its own when receives an `add` message.

[solution](counter/crdt/main.go) / [tests](counter/crdt/main_test.go)

Links:

- [CRDTs: The Hard Parts](https://www.youtube.com/watch?v=x7drE24geUw) by Martin Kleppmann - Great explanation about the differences between OTs and CRDTs
- [Conflict-free Replicated Data Types](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type) on Wikipedia

### [Challenge #5a: Single-Node Kafka-Style Log](https://fly.io/dist-sys/5a)

Simple implementation of a Kafka-style log. Store messages in a map, and return the messages after a given offset.

[solution](kafka/a/main.go) / [tests](kafka/a/main_test.go)

### [Challenge #5b: Multi-Node Kafka-Style Log](https://fly.io/dist-sys/5b)

Implementation using the Lin-KV store:

- Compare and Swap to get the next offset for a key
- Store the each message in a separated key
- Last Write Wins to store the committed offsets

[solution](kafka/b/main.go) / [tests](kafka/b/main_test.go)

### [Challenge #5c: Efficient Kafka-Style Log](https://fly.io/dist-sys/5c)

The previous solution also works for this challenge, so I tried to remove the KV completely and use a leader-follower approach.

The leader node is fixed to `n0` (could be done with a consensus algorithm).

The leader handles the offsets to avoid conflicts. All the nodes exchange messages to get updated messages and committed offsets.

[solution](kafka/c/main.go) / [tests](kafka/c/main_test.go)

## Project Structure

This repo is a Go workspace where each solution is in a separate module.

Each solution is in a folder with the `main.go` and `main_test.go` for tests.

### Tests

All solutions have automated tests using [custom utils](testutils) as a way to verify and play with the solutions, and to learn how to write tests in Go.

### Running

This project uses [Nix](nixos.org) for packages and [Taskfile](https://taskfile.dev) to run and build the solutions.

To run the solutions, use the `run-*` tasks:

```bash
task run-all
task run-uniqueids
task run-broadcast-efficient-ii1 # Will run and check the result output
```

To run the tests, use the `test-*` tasks:

```bash
task test-all
task test-uniqueids
task test-broadcast-efficient-ii
```
