# Test Circom

Benchmark circom with a maximum number of 8M constraints.

## Benchmark

- Intel(R) Xeon(R) E-2288G CPU @ 3.70GHz amd x86_64 Ubuntu 22.04 8 cores 128G Memory

## Usages

Install cicrom (refer to the official website) and snarkjs:

``` sh
npm install -g snarkjs
```

In the terminal:

``` sh
# Get Power-of-Tau (pot).
# These can be omitted for benchmark.
# snarkjs powersoftau new bn128 23 pot23_0000.ptau -v
# snarkjs powersoftau contribute pot23_0000.ptau pot23_0001.ptau --name="First contribution" -v
# snarkjs powersoftau beacon pot23_0001.ptau pot23_beacon.ptau 0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f 10 -n="Final Beacon"
# snarkjs powersoftau prepare phase2 pot23_beacon.ptau pot23.ptau -v
# We can borrow an existing pot:
# (./)
wget https://hermez.s3-eu-west-1.amazonaws.com/powersOfTau28_hez_final_23.ptau -O ./pot23.ptau

cd ./test_poseidon
```

Use ``nodejs`` and run the script,

``` sh
# (./test_poseidon)
bash ./bench.sh
```