# JS Witness Generation

## Install Dependencies

``` sh
npm install
```

## Usage

Suppose `main.circom` is the circuit for a specified user application.

1. (Offline) User generates circom files (`main.r1cs` and `./main_js/main.wasm`) via `circom` command line. 
    ``` sh
    circom main.circom --r1cs --json --wasm --sym --c --O0 
    ```

2. (Offline) User posts `main.r1cs` and `main.wasm` to Vortex Hub.

3. (Offline) Vortex Hub passes the two files to a Vortex Node. Node runs `Setup` for `main.r1cs` in `gnark`.

4. (Online) Client posts `input.json` to a Vortex Node. 
    - Node runs `witness_generator.js` to generate `main.wtns` file. Inputs: `main.r1cs`, `input.json` and `main.wasm`.
        ``` sh
        node witness_generator.js 
        ```
    - Node runs `Prove` for `main.r1cs` and `main.wtns` in `gnark` and generates a proof `proof`.
    - Node returns `proof` to client.
