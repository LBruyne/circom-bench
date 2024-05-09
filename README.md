# Circom Benchmark

Benchmark circom with Poseidon hashing function.
The maximum number of constraints to bench is 8M.

## Benchmark

### Server
- Circom version: 2.1.9
- Ubuntu 22.04 AMD x86_64 
- CPU: Intel(R) Xeon(R) E-2288G @ 3.70GHz 
- CPU cores: 8
- Memory: 128G

### PC

## Usages

Install cicrom (refer to the official website) and snarkjs:

``` sh
npm install -g snarkjs
```

For Groth16 proving purpose, the Circom library first need a ceremony called Power-of-Tau (pot/ptau). This can be done through the following commented steps.

However, the time cost of this setup is burdensome, therefore we borrow an existing pot from Internet (only for bench purpose).

In the terminal:

``` sh
# Get Power-of-Tau.
# Below is the ceremony, which can be omitted for benchmark.
# snarkjs powersoftau new bn128 23 pot23_0000.ptau -v
# snarkjs powersoftau contribute pot23_0000.ptau pot23_0001.ptau --name="First contribution" -v
# snarkjs powersoftau beacon pot23_0001.ptau pot23_beacon.ptau 0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f 10 -n="Final Beacon"
# snarkjs powersoftau prepare phase2 pot23_beacon.ptau pot23.ptau -v
#
# Instead, we can borrow an existing pot (this may be time cost):
# (./)
wget https://hermez.s3-eu-west-1.amazonaws.com/powersOfTau28_hez_final_23.ptau -O ./pot23.ptau

cd ./test_poseidon
```

The downloading may be time costing. To run a simple benchmark, Just modify the website url to adjust the number of tau smaller, (e.g., 14). However, this will restrict the maximum number of constraints we can bench.

Now we switch to ``nodejs`` and run the benchmark script,

``` sh
# (./test_poseidon)
bash ./bench.sh
```