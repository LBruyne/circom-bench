const snarkjs = require('snarkjs');
const { performance } = require('perf_hooks');

async function poseidon_test (hash_mode, hash_times, test_times) {
    const r1cs_file = `./${hash_times}/poseidon_${hash_mode}_1.r1cs`
    const wasm_file = `./${hash_times}/poseidon_${hash_mode}_1_js/poseidon_${hash_mode}_1.wasm`
    const ptau_file = "../pot22.ptau"
    // const prove_key = "./circuit.zkey"
    const input = {
        "inputs": Array.from({ length: hash_mode }, (_, index) => index + 1)
    }

    console.log(`Start benchmarking ${hash_times} "${hash_mode}-1 Poseidon Hash" and repeating ${test_times} times...`);

    // Setup
    const zkey = {type: "mem"}; 
    console.log("Groth16 setup...");
    await snarkjs.zKey.newZKey(r1cs_file, ptau_file, zkey);

    // Compute witness
    let start, duration_1, duration_2, duration_3;
    const wtns = {type: "mem"};
    console.log("Computing witness...");
    start = performance.now();
    for (let i = 0; i < test_times; i++) {
        await snarkjs.wtns.calculate(input, wasm_file, wtns);
    }
    duration_1 = performance.now() - start;
    console.log(`Total computing witness time: ${duration_1} ms`);

    // Prove 
    console.log("Groth16 proving...");
    start = performance.now();
    for (let i = 0; i < test_times; i++) {
        await snarkjs.groth16.prove(zkey, wtns);
    }
    duration_2 = performance.now() - start;
    console.log(`Total prove time: ${duration_2} ms`);

    // Full prove
    console.log("Groth16 full proving...");
    start = performance.now();
    for (let i = 0; i < test_times; i++) {
        const { proof, publicSignals } = await snarkjs.groth16.fullProve(
            input, 
            wasm_file, 
            zkey
        );
    }
    duration_3 = performance.now() - start;
    console.log(`Total fullProve time: ${duration_3} ms`);
    console.log(`Times: ${duration_1 / test_times} ms, ${duration_2 / test_times} ms, ${duration_3 / test_times}`);
    // console.log(publicSignals);
    // console.log(proof);

    console.log("End.");
}

(async function main() {
    // Checking if command line arguments are provided
    if (process.argv.length < 4) {
        console.log("Usage: node <scriptname> <hash_mode> <hash_times> <test_times>");
        process.exit(1);
    }
    
    const hash_mode = parseInt(process.argv[2], 10); // hash_mode is the mode of Poseidon hash, meaning the number of inputs. Must be 2/4/8/16.
    const hash_times = parseInt(process.argv[3], 10);
    const test_times = parseInt(process.argv[4], 10);

    await poseidon_test(hash_mode, hash_times, test_times);
    process.exit(0);
})()
