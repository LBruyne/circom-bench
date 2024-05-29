#!/bin/bash

CIRCUIT_NAMES=("sudoku" "rollup" "poseidon_16" "eth_addr" "ecdsa_verify")

for CIRCUIT_NAME in "${CIRCUIT_NAMES[@]}"; do
    echo "****COMPILING CIRCUIT $CIRCUIT_NAME****"
    start=$(date +%s)
    circom "$CIRCUIT_NAME.circom" --r1cs --wasm --sym --c --O0
    end=$(date +%s)
    echo "Total compile circom time: ($((end-start))s)"

    echo "****GENERATE WITNESS FOR $CIRCUIT_NAME****"
    node witness_generator.js "$CIRCUIT_NAME"
    echo "Done."
done
