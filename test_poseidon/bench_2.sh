#!/bin/bash

hash_mode=16
hash_times=(300 600 900 1000)

for folder in "${hash_times[@]}"; do

    cd "$folder" || exit

    echo "Start benchmarking $folder ""$hash_mode"-1 Poseidon Hash"..."

    # Setup
    snarkjs groth16 setup ./poseidon_16_1.r1cs ../../pot22.ptau circuit.zkey 

    # Computing witness
    start_1=$(date +%s%N) 
    node ./poseidon_16_1_js/generate_witness.js ./poseidon_16_1_js/poseidon_16_1.wasm ../input.json witness.wtns
    end_1=$(date +%s%N) 

    duration_1=$(( (end_1 - start_1) / 1000000 ))

    echo "Compute witness time: $duration_1 ms"

    # Proving
    start_2=$(date +%s%N) 
    snarkjs groth16 prove ./circuit.zkey witness.wtns proof.json public.json
    end_2=$(date +%s%N) 

    duration_2=$(( (end_2 - start_2) / 1000000 ))

    echo "Prove time: $duration_2 ms"

    # Full proving
    start_3=$(date +%s%N) 
    snarkjs groth16 fullprove ../input.json ./poseidon_16_1_js/poseidon_16_1.wasm ./circuit.zkey proof.json public.json
    end_3=$(date +%s%N) 

    duration_3=$(( (end_3 - start_3) / 1000000 ))

    echo "Full prove time: $duration_3 ms"

    echo "End."

    cd ..
done