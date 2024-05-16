#!/bin/bash

#hash_times=(3 10 30 100 150 200 230 250 300 600 900 1000)
hash_times=(1000)

for folder in "${hash_times[@]}"; do
    cd "$folder" || exit

    if [[ -f "poseidon_16_1.circom" ]]; then
        circom poseidon_16_1.circom --r1cs --json --wasm --sym --c --O0 
    else
        echo "File poseidon_16_1.circom not found in folder $folder"
    fi

    cd ..
done

