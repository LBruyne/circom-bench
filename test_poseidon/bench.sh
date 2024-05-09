#!/bin/bash

hash_times=(3 10 30 100 150 200 230 250) 
test_times=(10 10 10 4 4 4 4 4)

# hash_times=(300 600)
# test_times=(2 2)

if [ "${#hash_times[@]}" -ne "${#test_times[@]}" ]; then
    echo "Error: The length of hash_times and test_times arrays must be the same."
    exit 1
fi

for folder in "${hash_times[@]}"; do
    cd "$folder" || exit

    if [[ -f "poseidon_16_1.circom" ]]; then
        circom poseidon_16_1.circom --r1cs --json --wasm --sym --c --O0 
    else
        echo "File poseidon_16_1.circom not found in folder $folder"
    fi

    cd ..
done

for i in "${!hash_times[@]}"; do
    node ./run.js 16 ${hash_times[i]} ${test_times[i]}
done