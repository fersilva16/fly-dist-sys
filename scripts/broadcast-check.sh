#! /usr/bin/env bash

msgs_per_op=$(cat store/broadcast/latest/results.edn | jet -t ":net :servers :msgs-per-op")

if (( $(echo "$msgs_per_op > $1" | bc -l) )); then
  echo "Messages per op is greater than expected ($msgs_per_op/$1)"

  exit 1
fi

med_latency=$(cat store/broadcast/latest/results.edn | jet -T ":workload -> :stable-latencies -> (get 0.5)")

if (( $(echo "$med_latency > $2" | bc -l) )); then
  echo "Median latency is greater than expected ($med_latency/$2)"

  exit 1
fi

max_latency=$(cat store/broadcast/latest/results.edn | jet -T ":workload -> :stable-latencies -> (get 1)")

if (( $(echo "$max_latency > $3" | bc -l) )); then
  echo "Maximum latency is greater than expected ($max_latency/$3)"

  exit 1
fi

echo "All checks passed"
echo "  Messages per op: $msgs_per_op/$1"
echo "  Median latency: $med_latency/$2"
echo "  Maximum latency: $max_latency/$3"