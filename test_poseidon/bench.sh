#!/bin/bash

hash_times=(3 10 30 100 150 200 230 250) 
test_times=(10 10 10 4 4 4 4 4)

# Part 1
if [ "${#hash_times[@]}" -ne "${#test_times[@]}" ]; then
    echo "Error: The length of hash_times and test_times arrays must be the same."
    exit 1
fi

for i in "${!hash_times[@]}"; do
    node ./run.js 16 ${hash_times[i]} ${test_times[i]}
done

# Part 2
hash_mode=16
hash_times=(300 600 900 1000)

for folder in "${hash_times[@]}"; do

    cd "$folder" || exit

    echo "Start benchmarking $folder ""$hash_mode"-1 Poseidon Hash"..."

    # Setup
    snarkjs groth16 setup ./poseidon_16_1.r1cs ../../pot22.ptau circuit.zkey 

    # Computing witness
    start_1=$(date +%s) 
    node ./poseidon_16_1_js/generate_witness.js ./poseidon_16_1_js/poseidon_16_1.wasm ../input.json witness.wtns
    end_1=$(date +%s) 

    echo "Start: $start_1"
    echo "End: $end_1"

    duration_1=$(( (end_1 - start_1) ))

    echo "Compute witness time: $duration_1 s"

    # Proving
    start_2=$(date +%s) 
    snarkjs groth16 prove ./circuit.zkey witness.wtns proof.json public.json
    end_2=$(date +%s) 

    duration_2=$(( (end_2 - start_2) ))

    echo "Prove time: $duration_2 s"

    # Full proving
    start_3=$(date +%s) 
    snarkjs groth16 fullprove ../input.json ./poseidon_16_1_js/poseidon_16_1.wasm ./circuit.zkey proof.json public.json
    end_3=$(date +%s) 

    duration_3=$(( (end_3 - start_3) ))

    echo "Full prove time: $duration_3 s"

    echo "End."

    cd ..
done