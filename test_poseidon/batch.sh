#!/bin/bash

parallel=$1

for ((i = 0; i < parallel; i++)); do
    ./batch-bench.sh $i &
done

for job in $(jobs -p); do
    wait $job
done
