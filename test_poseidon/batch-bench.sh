#!/bin/bash

index=$1

hash_mode=16
hash_times=(900 1000)

for folder in "${hash_times[@]}"; do

    cd "$folder" || exit

    echo "Start benchmarking $folder ""$hash_mode"-1 Poseidon Hash"..."

    # Setup
    potname="../../pot22-$index.ptau"
    keyname="circuit-$index.zkey"
    snarkjs groth16 setup ./poseidon_16_1.r1cs $potname $keyname 

    # Full proving
    start_3=$(date +%s) 
    inputname="../input-$index.json"
    snarkjs groth16 fullprove $inputname ./poseidon_16_1_js/poseidon_16_1.wasm $keyname proof.json public.json
    end_3=$(date +%s) 

    duration_3=$(( (end_3 - start_3) ))

    echo "Full prove time: $duration_3 s"

    echo "End."

    cd ..
done
