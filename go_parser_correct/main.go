package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	bn254_system "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"go_parser_correct/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"time"
)

func main() {
	//If u want to generate witness with witness_generator , u can use the function;before u use it
	//u should finish the prerequisites for js_witness in readme ↓↓↓↓
	//RunGenerateWnts()

	//U can run the function Example to to complete the proof process ↓↓↓↓ （这个函数目前存在问题，无法加载读取保存的json r1cs，有待解决）
	//Example()

	// U can perform multiple rounds of testing using this function for one prove mission ↓↓↓↓
	// —— For multi files
	//basePath := "../js_witness/"
	//fileNames := []string{"rollup", "sudoku", "poseidon_16", "eth_addr", "ecdsa_verify"}
	//for _, fileName := range fileNames {
	//	fmt.Println("\nTest [", fileName, "]")
	//	testFunction(filepath.Join(basePath, fileName+".r1cs"), filepath.Join(basePath, fileName+".wtns"), 5)
	//}

	testFunction("../js_witness/sudoku.r1cs", "../js_witness/sudoku.wtns", 5)
	testFunction("../js_witness/rollup.r1cs", "../js_witness/rollup.wtns", 5)
	testFunction("../js_witness/eth_addr.r1cs", "../js_witness/eth_addr.wtns", 5)
	testFunction("../js_witness/ecdsa_verify.r1cs", "../js_witness/ecdsa_verify.wtns", 5)
	testFunction("../js_witness/poseidon_16.r1cs", "../js_witness/poseidon_16.wtns", 5)

	// —— For single file
	//testFunction("../js_witness/main_c.r1cs", "../js_witness/main_c.wtns", 5)
}

func RunGenerateWnts() {
	scriptPath, err := filepath.Abs("../js_witness/witness_generator.js")
	if err != nil {
		panic(err)
	}
	scriptDir := filepath.Dir(scriptPath)
	cmd := exec.Command("node", scriptPath)
	cmd.Dir = scriptDir
	//cmd := exec.Command("node", "../test_poseidon/generate_witness.js", "../test_poseidon/poseidon_16_1.wasm", "../test_poseidon/input.json", "./utils/output.wtns")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", output)
}
func Example() {
	startRead := time.Now()
	r1cs, err := ReadR1CS("../js_witness/main_c.r1cs")
	durationRead := time.Since(startRead)
	fmt.Printf("Read time: %v\n", durationRead)
	if err != nil {
		panic(err)
	}
	var w R1CSCircuit
	startParse := time.Now()
	w.Witness, err = utils.ParseWtns("../js_witness/main_c.wtns")
	durationParse := time.Since(startParse)
	fmt.Printf("Parse time: %v\n", durationParse)
	if err != nil {
		panic(err)
	}
	secretWitness, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}
	witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		panic(err)
	}
	startSetup := time.Now()
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		panic(err)
	}
	durationSetup := time.Since(startSetup)
	fmt.Printf("Setup time: %v\n", durationSetup)
	SaveTime := time.Now()
	err = SaveToJSON("Pk.json", pk)
	err = SaveToJSON("Vk.json", vk)
	var buffer bytes.Buffer
	r1cs.WriteTo(&buffer)
	err = ioutil.WriteFile("r1cs_data.bin", buffer.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing buffer to file:", err)
		return
	}

	data, err := ioutil.ReadFile("r1cs_data.bin")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	buffer = *bytes.NewBuffer(data)
	var reconstructedR1cs bn254_system.R1CS
	_, err = reconstructedR1cs.ReadFrom(&buffer)
	if err != nil {
		panic(err)
	}
	durationSave := time.Since(SaveTime)
	fmt.Printf("Save Pk Vk R1cs time: %v\n", durationSave)
	var deserializedPK groth16_bn254.ProvingKey
	var deserializedVK groth16_bn254.VerifyingKey
	var deserializedProof groth16_bn254.Proof
	ReadJsonTime := time.Now()
	err = LoadFromJSON("Pk.json", &deserializedPK)
	if err != nil {
		fmt.Println("Error deserializing:", err)
		return
	}
	durationReadJson := time.Since((ReadJsonTime))
	fmt.Printf("ReadPKJson time: %v\n", durationReadJson)
	ReadJsonTime = time.Now()
	err = LoadFromJSON("Vk.json", &deserializedVK)
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}
	startProve := time.Now()
	proof, err := groth16.Prove(&reconstructedR1cs, &deserializedPK, secretWitness, backend.WithIcicleAcceleration())
	err = SaveToJSON("proof.json", proof)
	if err != nil {
		panic(err)
	}
	err = LoadFromJSON("proof.json", &deserializedProof)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	durationProve := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", durationProve)
	startVerify := time.Now()
	_ = groth16.Verify(&deserializedProof, &deserializedVK, witnessPublic)
	durationVerify := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", durationVerify)
}

func testFunction(r1cs_path string, wtns_path string, prove_cycles int) {
	ccs, err := ReadR1CS(r1cs_path)
	if err != nil {
		panic(err)
	}
	var w R1CSCircuit
	w.Witness, err = utils.ParseWtns(wtns_path)
	if err != nil {
		panic(err)
	}
	secretWitness, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}
	witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		panic(err)
	}
	startSetup := time.Now()
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		panic(err)
	}
	durationSetup := time.Since(startSetup)
	fmt.Printf("Setup time: %v\n", durationSetup)

	var totalTime int64
	for i := 0; i < prove_cycles; i++ {
		fmt.Printf("——————————prove time %v ——————————\n", i+1)
		startProve := time.Now()
		proof, err := groth16.Prove(ccs, pk, secretWitness)
		if err != nil {
			panic(err)
		}
		durationProve := time.Since(startProve)
		fmt.Printf("Prove time: %v\n", durationProve)
		startVerify := time.Now()
		err = groth16.Verify(proof, vk, witnessPublic)
		durationVerify := time.Since(startVerify)
		fmt.Printf("Verify time: %v\n", durationVerify)
		fmt.Printf("—————————————————————————————————\n")
		totalTime += int64(durationProve)
	}
	fmt.Printf("Average prove time : %v\n", time.Duration(totalTime/int64(prove_cycles)))
}

func SaveToJSON(filePath string, v interface{}) error {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadFromJSON(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("v must be a non-nil pointer")
	}
	if err = json.Unmarshal(byteValue, v); err != nil {
		return err
	}
	return nil
}
