# Circom Benchmark

Benchmark circom and snarkjs with Poseidon hashing function.
The maximum number of constraints to bench is 4M.

## Benchmark

### In Server \#1

Setting:
- Circom version: 2.1.9
- SnarkJS version 0.7.4
- Ubuntu 22.04 AMD x86_64 
- CPU: Intel(R) Xeon(R) E-2288G @ 3.70GHz 
- CPU cores: 8
- Threads: 16
- Memory: 128G

Note: SnarkJS does enable multi-thread optimization.

### In Server \#2

Setting:
- Circom version 2.1.9
- SnarkJS version 0.7.4
- Ubuntu 22.04 AMD x86_64 
- CPU: Intel(R) Xeon(R) Gold 6462C @ ? GHz
- CPU cores: 32
- Threads: 64
- Memory: 252G

### In PC

Setting:
- Circom version 2.1.2
- SnarkJS version 0.7.4
- MacOS Sonoma 14.3.1
- CPU: Apple M1 Pro
- CPU cores: 10
- Threads: 10
- Memory: 16G

## Usages

Install circom (refer to the official website) and snarkjs:

``` sh
npm install -g snarkjs
```

For Groth16 proving purpose, the Circom library first need a ceremony called Power-of-Tau (pot/ptau). This can be done through the following commented steps.
However, the time cost of this setup is burdensome, therefore we borrow an existing pot from Internet (only for bench purpose).

In the terminal:

``` sh
# Get Power-of-Tau.
# Below is the ceremony, which can be omitted for benchmark.
# snarkjs powersoftau new bn128 22 pot22_0000.ptau -v
# snarkjs powersoftau contribute pot23_0000.ptau pot22_0001.ptau --name="First contribution" -v
# snarkjs powersoftau beacon pot22_0001.ptau pot22_beacon.ptau 0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f 10 -n="Final Beacon"
# snarkjs powersoftau prepare phase2 pot22_beacon.ptau pot22.ptau -v
#
# Instead, we can borrow an existing pot (this may be time cost):
# (./)
wget https://hermez.s3-eu-west-1.amazonaws.com/powersOfTau28_hez_final_22.ptau -O ./pot22.ptau

cd ./test_poseidon
```

The downloading may be time-costing. To run a simple benchmark, just modify the website url to adjust the number of tau smaller (e.g., 16). However, this restricts the maximum number of constraints we can bench.

Now switch to ``nodejs`` and run the benchmark script,

``` sh
# (./test_poseidon)
# Prepare circuits
bash ./prepare.sh

# (./test_poseidon)
# This benchmark 3-1000 Poseidon hash (with maximum 4000000 constraints)
bash ./bench.sh
```
