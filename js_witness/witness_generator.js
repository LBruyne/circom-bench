const snarkjs = require('snarkjs');
const fs = require('fs');
const { performance } = require('perf_hooks');

async function witness_generator () {
    const r1cs_file = './main.r1cs'
    const wasm_file = './main_js/main.wasm'
    const wtns_file = './main.wtns'
    const input_file = './input.json'
    let start, duration

    console.log(`Start generating witness for ${r1cs_file} circuit using ${input_file} as client inputs, with the aid of ${wasm_file} WASM file.`);

    // Load input
    console.log(`Load ${input_file}.`);
    start = performance.now();
    const input = JSON.parse(fs.readFileSync(input_file, 'utf8'))
    duration = performance.now() - start;
    console.log(`Total loading input time: ${duration} ms`);

    // Compute witness
    console.log("Computing witness.");
    start = performance.now();
    await snarkjs.wtns.calculate(input, wasm_file, wtns_file);
    duration = performance.now() - start;
    console.log(`Total computing witness time: ${duration} ms`);
    console.log("End.");
}

(async function main() {
    await witness_generator();
    process.exit(0);
})()
