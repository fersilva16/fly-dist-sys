#! /usr/bin/env bash

msgs_per_op=$(cat store/broadcast/latest/results.edn | jet -t ":net :servers :msgs-per-op")

if (( $(echo "$msgs_per_op > $1" | bc -l) )); then
  echo "Messages per op is greater than expected ($1/$msgs_per_op)"

  exit 1
fi

med_latency=$(cat store/broadcast/latest/results.edn | jet -T ":workload -> :stable-latencies -> (get 0.5)")

if (( $(echo "$med_latency > $2" | bc -l) )); then
  echo "Median latency is greater than expected ($2/$med_latency)"

  exit 1
fi

max_latency=$(cat store/broadcast/latest/results.edn | jet -T ":workload -> :stable-latencies -> (get 1)")

if (( $(echo "$max_latency > $3" | bc -l) )); then
  echo "Maximum latency is greater than expected ($3/$max_latency)"

  exit 1
fi