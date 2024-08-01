package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	fr_bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"go_parser_correct/utils"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strings"
	"time"
)

//	func main2() {
//		fileChannel := make(chan string)
//		done := make(chan bool)
//		exit := make(chan os.Signal, 1)
//		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
//		r1csPath := "../js_witness/poseidon_16.r1cs" // 替换为对应的 R1CS 文件路径
//		ccs, err := ReadR1CS(r1csPath)
//		if err != nil {
//			log.Fatalf("Failed to read R1CS file: %v", err)
//		}
//		startSetup := time.Now()
//		pk, vk, err := groth16.Setup(ccs)
//		if err != nil {
//			log.Fatalf("Failed to setup pk and vk: %v", err)
//		}
//		durationSetup := time.Since(startSetup)
//		fmt.Printf("Setup time: %v\n", durationSetup)
//
//		// 启动一个 goroutine 处理文件
//		go func() {
//			for {
//				select {
//				case wtnsPath := <-fileChannel:
//					RunProve_con(ccs, pk, vk, wtnsPath)
//					done <- true
//				case <-exit:
//					fmt.Println("Shutting down...")
//					return
//				}
//			}
//		}()
//		go func() {
//			reader := bufio.NewReader(os.Stdin)
//			for {
//				fmt.Printf("—————————————————————————————————\n")
//				fmt.Print("Enter wtns file path: ")
//				wtnsPath, _ := reader.ReadString('\n')
//				wtnsPath = strings.TrimSpace(wtnsPath)
//				if wtnsPath == "exit" {
//					close(exit)
//					return
//				}
//				fileChannel <- wtnsPath
//				success := <-done
//				if success {
//					fmt.Printf("File %s processed successfully.\n", wtnsPath)
//				} else {
//					fmt.Printf("Failed to process file %s.\n", wtnsPath)
//				}
//			}
//		}()
//
//		<-exit
//		fmt.Println("Program terminated.")
//	}
func main() {
	//If u want to generate witness with witness_generator , u can use the function;before u use it
	//u should finish the prerequisites for js_witness in readme ↓↓↓↓
	//RunGenerateWnts()

	//U can run the function Example to complete the proof process ↓↓↓↓ （这个函数目前存在问题，无法加载读取保存的json r1cs，有待解决，已解决）
	//Example()

	// U can perform multiple rounds of testing using this function for one prove mission ↓↓↓↓
	// —— For multi files
	//basePath := "../js_witness/"
	//fileNames := []string{"rollup", "sudoku", "poseidon_16", "eth_addr", "ecdsa_verify"}
	//for _, fileName := range fileNames {
	//	fmt.Println("\nTest [", fileName, "]")
	//	testFunction(filepath.Join(basePath, fileName+".r1cs"), filepath.Join(basePath, fileName+".wtns"), 5)
	//}

	//testFunction("../js_witness/sudoku.r1cs", "../js_witness/sudoku.wtns", 1)
	//testFunction("../js_witness/rollup.r1cs", "../js_witness/rollup.wtns", 1)
	//testFunction("../js_witness/rollup.r1cs", "../js_witness/rollup.wtns", 1)
	//testFunction("../js_witness/eth_addr.r1cs", "../js_witness/eth_addr.wtns", 1)
	//testFunction("../js_witness/ecdsa_verify.r1cs", "../js_witness/ecdsa_verify.wtns", 5)
	//testFunction("../test_poseidon/10/poseidon_16_1.r1cs", "../js_witness/poseidon_16.wtns", 1)
	//RunProve("../test_poseidon/10/poseidon_12_20.r1cs", "../test_poseidon/10/poseidon_12_20_js/output.wtns")
	RunProve("./ZKLogin.r1cs", "./witness.wtns")
	//RunProve_2("../test_poseidon/10/poseidon_16_1.r1cs", "../test_poseidon/10/poseidon_16_1_js/output3.wtns", "../test_poseidon/10/poseidon_16_1_js/output2.wtns")
	//RunProve("../js_witness/rollup.r1cs", "../js_witness/rollup.wtns")
	//Example()
	// —— For single file
	//testFunction("../js_witness/main_c.r1cs", "../js_witness/main_c.wtns", 5)
}

//	func RunGenerateWnts() {
//		scriptPath, err := filepath.Abs("../js_witness/witness_generator.js")
//		if err != nil {
//			panic(err)
//		}
//		scriptDir := filepath.Dir(scriptPath)
//		cmd := exec.Command("node", scriptPath)
//		cmd.Dir = scriptDir
//		//cmd := exec.Command("node", "../test_poseidon/generate_witness.js", "../test_poseidon/poseidon_16_1.wasm", "../test_poseidon/input.json", "./utils/output.wtns")
//		output, err := cmd.CombinedOutput()
//		if err != nil {
//			fmt.Printf("Error: %s\n", err)
//			return
//		}
//		fmt.Printf("Output: %s\n", output)
//	}
//
//	func Example() {
//		startRead := time.Now()
//		r1cs, err := ReadR1CS("../js_witness/main_c.r1cs")
//		durationRead := time.Since(startRead)
//		fmt.Printf("Read time: %v\n", durationRead)
//		if err != nil {
//			panic(err)
//		}
//		var w R1CSCircuit
//		startParse := time.Now()
//		w.Witness, err = utils.ParseWtns("../js_witness/main_c.wtns")
//		durationParse := time.Since(startParse)
//		fmt.Printf("Parse time: %v\n", durationParse)
//		if err != nil {
//			panic(err)
//		}
//		secretWitness, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
//		if err != nil {
//			panic(err)
//		}
//		witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
//		if err != nil {
//			panic(err)
//		}
//		startSetup := time.Now()
//		pk, vk, err := groth16.Setup(r1cs)
//		if err != nil {
//			panic(err)
//		}
//		durationSetup := time.Since(startSetup)
//		fmt.Printf("Setup time: %v\n", durationSetup)
//		SaveTime := time.Now()
//		err = SaveToJSON("Pk.json", pk)
//		err = SaveToJSON("Vk.json", vk)
//		var buffer bytes.Buffer
//		r1cs.WriteTo(&buffer)
//		err = ioutil.WriteFile("r1cs_data.bin", buffer.Bytes(), 0644)
//		if err != nil {
//			fmt.Println("Error writing buffer to file:", err)
//			return
//		}
//
//		data, err := ioutil.ReadFile("r1cs_data.bin")
//		if err != nil {
//			fmt.Println("Error reading file:", err)
//			return
//		}
//		buffer = *bytes.NewBuffer(data)
//		var reconstructedR1cs bn254_system.R1CS
//		_, err = reconstructedR1cs.ReadFrom(&buffer)
//		if err != nil {
//			panic(err)
//		}
//		durationSave := time.Since(SaveTime)
//		fmt.Printf("Save Pk Vk R1cs time: %v\n", durationSave)
//		var deserializedPK groth16_bn254.ProvingKey
//		var deserializedVK groth16_bn254.VerifyingKey
//		var deserializedProof groth16_bn254.Proof
//		ReadJsonTime := time.Now()
//		err = LoadFromJSON("Pk.json", &deserializedPK)
//		if err != nil {
//			fmt.Println("Error deserializing:", err)
//			return
//		}
//		durationReadJson := time.Since((ReadJsonTime))
//		fmt.Printf("ReadPKJson time: %v\n", durationReadJson)
//		ReadJsonTime = time.Now()
//		err = LoadFromJSON("Vk.json", &deserializedVK)
//		if err != nil {
//			fmt.Println("Error reading from file:", err)
//			return
//		}
//		startProve := time.Now()
//		proof, err := groth16.Prove(&reconstructedR1cs, &deserializedPK, secretWitness, backend.WithIcicleAcceleration())
//		if p, ok := proof.(*groth16_bn254.Proof); ok {
//			// 访问字段
//			fmt.Println("Ar:", p.Ar)
//			fmt.Println("Krs:", p.Krs)
//			fmt.Println("Bs:", p.Bs)
//			fmt.Println("Commitments:", p.Commitments)
//			fmt.Println("CommitmentPok:", p.CommitmentPok)
//		} else {
//			fmt.Println("类型断言失败")
//		}
//		err = SaveToJSON("proof.json", proof)
//		if err != nil {
//			panic(err)
//		}
//		err = LoadFromJSON("proof.json", &deserializedProof)
//		if err != nil {
//			panic(err)
//		}
//		if err != nil {
//			panic(err)
//		}
//		durationProve := time.Since(startProve)
//		fmt.Printf("Prove time: %v\n", durationProve)
//		startVerify := time.Now()
//		_ = groth16.Verify(&deserializedProof, &deserializedVK, witnessPublic)
//		durationVerify := time.Since(startVerify)
//		fmt.Printf("Verify time: %v\n", durationVerify)
//	}
//
//	func testFunction(r1cs_path string, wtns_path string, prove_cycles int) {
//		ccs, err := ReadR1CS(r1cs_path)
//		if err != nil {
//			panic(err)
//		}
//		var w R1CSCircuit
//		w.Witness, err = utils.ParseWtns(wtns_path)
//		if err != nil {
//			panic(err)
//		}
//		secretWitness, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
//		if err != nil {
//			panic(err)
//		}
//		witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
//		if err != nil {
//			panic(err)
//		}
//		startSetup := time.Now()
//		pk, vk, err := groth16.Setup(ccs)
//		if err != nil {
//			panic(err)
//		}
//		durationSetup := time.Since(startSetup)
//		fmt.Printf("Setup time: %v\n", durationSetup)
//
//		var totalTime int64
//		for i := 0; i < prove_cycles; i++ {
//			fmt.Printf("——————————prove time %v ——————————\n", i+1)
//			startProve := time.Now()
//			proof, err := groth16.Prove(ccs, pk, secretWitness)
//			if err != nil {
//				panic(err)
//			}
//			durationProve := time.Since(startProve)
//			fmt.Printf("Prove time: %v\n", durationProve)
//			startVerify := time.Now()
//			err = groth16.Verify(proof, vk, witnessPublic)
//			durationVerify := time.Since(startVerify)
//			fmt.Printf("Verify time: %v\n", durationVerify)
//			fmt.Printf("—————————————————————————————————\n")
//			totalTime += int64(durationProve)
//		}
//		fmt.Printf("Average prove time : %v\n", time.Duration(totalTime/int64(prove_cycles)))
//	}
func RunProve(r1cs_path string, wtns_path string) {
	ccs, NumOutput, NumInPublic, err := ReadR1CS(r1cs_path)
	if err != nil {
		panic(err)
	}
	var w R1CSCircuit
	w.Witness, w.WitnessPublic, err = utils.ParseWtns(wtns_path, NumOutput, NumInPublic)
	if err != nil {
		panic(err)
	}

	//compare
	//var w2 R1CSCircuit
	//w2.Witness = w.Witness
	//w2.WitnessPublic = make([]frontend.Variable, len(w.WitnessPublic))
	//copy(w2.WitnessPublic, w.WitnessPublic)
	//w2.WitnessPublic[5] = big.NewInt(250)
	//witnessPublic2, err := frontend.NewWitness(&w2, ecc.BN254.ScalarField(), frontend.PublicOnly())
	//if err != nil {
	//	panic(err)
	//}

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
	startProve := time.Now()
	proof, err := groth16.Prove(ccs, pk, secretWitness)
	if err != nil {
		panic(err)
	}
	durationProve := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", durationProve)
	//err = SaveToJSON("proof.json", proof)
	//if err != nil {
	//	panic(err)
	//}
	startVerify := time.Now()
	//var deserializedProof groth16_bn254.Proof
	//err = LoadFromJSON("proof.json", &deserializedProof)
	//if err != nil {
	//	panic(err)
	//}
	err = groth16.Verify(proof, vk, witnessPublic)
	if err != nil {
		panic(err)
	}

	//err = groth16.Verify(proof, vk, witnessPublic2)
	//if err != nil {
	//	fmt.Println("verify fail:" + err.Error())
	//}
	//err = groth16.Verify(&deserializedProof, vk, witnessPublic)
	//if err != nil {
	//	panic(err)
	//}
	durationVerify := time.Since(startVerify)
	f, err := os.Create("verifier.sol")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := vk.ExportSolidity(f); err != nil {
		panic(err)
	}
	fmt.Printf("Verify time: %v\n", durationVerify)
	fmt.Printf("—————————————————————————————————\n")
	totalTime += int64(durationProve)
	fmt.Printf("Average prove time : %v\n", time.Duration(totalTime))

	//for i := 0; i < 5; i++ {
	//	startProve = time.Now()
	//	_, err := groth16.Prove(ccs, pk, secretWitness)
	//	if err != nil {
	//		panic(err)
	//	}
	//	durationProve = time.Since(startProve)
	//	fmt.Printf("Prove time: %v\n", durationProve)
	//}

	// proof to hex
	_proof, ok := proof.(interface{ MarshalSolidity() []byte })
	if !ok {
		panic("proof does not implement MarshalSolidity()")
	}
	proofStr := hex.EncodeToString(_proof.MarshalSolidity())

	// public witness to hex
	bPublicWitness, err := witnessPublic.MarshalBinary()
	if err != nil {
		panic(err)
	}
	// that's quite dirty...
	// first 4 bytes -> nbPublic
	// next 4 bytes -> nbSecret
	// next 4 bytes -> nb elements in the vector (== nbPublic + nbSecret)
	bPublicWitness = bPublicWitness[12:]
	publicWitnessStr := hex.EncodeToString(bPublicWitness)

	// https://github.com/Consensys/gnark-solidity-checker/blob/main/cmd/templates.go

	proofHex := proofStr
	inputHex := publicWitnessStr
	nbPublicInputs := len(witnessPublic.Vector().(fr_bn254.Vector))
	fpSize := 4 * 8

	proofBytes, err := hex.DecodeString(proofHex)
	if err != nil {
		panic(err)
	}

	if len(proofBytes) != fpSize*8 {
		panic("proofBytes != fpSize*8")
	}

	inputBytes, err := hex.DecodeString(inputHex)
	if err != nil {
		panic(err)
	}

	if len(inputBytes)%fr.Bytes != 0 {
		panic("inputBytes mod fr.Bytes !=0")
	}

	// convert public inputs
	nbInputs := len(inputBytes) / fr.Bytes
	if nbInputs != nbPublicInputs {
		panic("nbInputs != nbPublicInputs")
	}
	input := make([]*big.Int, nbPublicInputs)
	for i := 0; i < nbInputs; i++ {
		var e fr.Element
		e.SetBytes(inputBytes[fr.Bytes*i : fr.Bytes*(i+1)])
		input[i] = new(big.Int)
		e.BigInt(input[i])
	}

	// solidity contract inputs
	var __proof [8]*big.Int

	// proof.Ar, proof.Bs, proof.Krs
	for i := 0; i < 8; i++ {
		__proof[i] = new(big.Int).SetBytes(proofBytes[fpSize*i : fpSize*(i+1)])
	}

	// Create a string slice that stores each big.Int converted string
	var finalproof []string
	var finalinput []string

	for _, bi := range __proof {
		finalproof = append(finalproof, bi.String())
	}
	for _, bi := range input {
		finalinput = append(finalinput, bi.String())
	}
	result := "[" + strings.Join(finalproof, ",") + "]," + "[" + strings.Join(finalinput, ",") + "]"

	fmt.Println("result: ", result)
}

func RunProve_2(r1cs_path string, wtns_path string, wtns_path_2 string) {
	ccs, NumOutput, NumInPublic, err := ReadR1CS(r1cs_path)
	if err != nil {
		panic(err)
	}
	var w, w2 R1CSCircuit
	w.Witness, w.WitnessPublic, err = utils.ParseWtns(wtns_path, NumOutput, NumInPublic)
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
	w2.Witness, w2.WitnessPublic, err = utils.ParseWtns(wtns_path_2, NumOutput, NumInPublic)
	if err != nil {
		panic(err)
	}
	secretWitness2, err := frontend.NewWitness(&w2, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}
	witnessPublic2, err := frontend.NewWitness(&w2, ecc.BN254.ScalarField(), frontend.PublicOnly())
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

	startProve := time.Now()
	proof, err := groth16.Prove(ccs, pk, secretWitness)
	if err != nil {
		panic(err)
	}
	durationProve := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", durationProve)
	err = SaveToJSON("proof1.json", proof)
	if err != nil {
		panic(err)
	}
	startVerify := time.Now()
	err = groth16.Verify(proof, vk, witnessPublic)
	durationVerify := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", durationVerify)
	fmt.Printf("—————————————————————————————————\n")

	startProve = time.Now()
	proof2, err := groth16.Prove(ccs, pk, secretWitness2)
	if err != nil {
		panic(err)
	}
	durationProve = time.Since(startProve)
	fmt.Printf("Prove time: %v\n", durationProve)
	err = SaveToJSON("proof2.json", proof2)
	if err != nil {
		panic(err)
	}
	startVerify = time.Now()
	err = groth16.Verify(proof, vk, witnessPublic)
	if err != nil {
		panic(err)
	}
	//err = groth16.Verify(proof2, vk, witnessPublic)
	//if err != nil {
	//	panic(err)
	//}
	err = groth16.Verify(proof2, vk, witnessPublic2)
	if err != nil {
		panic(err)
	}
	durationVerify = time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", durationVerify)
	fmt.Printf("—————————————————————————————————\n")

	//for i := 0; i < 5; i++ {
	//	_, err = groth16.Prove(ccs, pk, secretWitness)
	//	if err != nil {
	//		panic(err)
	//	}
	//	_, err = groth16.Prove(ccs, pk, secretWitness2)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
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

//func RunProve_con(ccs constraint.ConstraintSystem, pk groth16.ProvingKey, vk groth16.VerifyingKey, wtnsPath string) {
//	var w R1CSCircuit
//	var err error
//	w.Witness, err = utils.ParseWtns(wtnsPath)
//	if err != nil {
//		panic(err)
//	}
//	secretWitness, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
//	if err != nil {
//		panic(err)
//	}
//	witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
//	if err != nil {
//		panic(err)
//	}
//	startProve := time.Now()
//	proof, err := groth16.Prove(ccs, pk, secretWitness)
//	if err != nil {
//		panic(err)
//	}
//	durationProve := time.Since(startProve)
//	fmt.Printf("Prove time: %v\n", durationProve)
//
//	baseName := filepath.Base(wtnsPath)
//	fileNameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
//	filename := fmt.Sprintf("%s_proof.json", fileNameWithoutExt)
//
//	//uuid := uuid.New().String()[:6]
//	//filename := fmt.Sprintf("proof_%s.json", uuid)
//	err = SaveToJSON(filename, proof)
//	if err != nil {
//		panic(err)
//	}
//	startVerify := time.Now()
//	err = groth16.Verify(proof, vk, witnessPublic)
//	durationVerify := time.Since(startVerify)
//	fmt.Printf("Verify time: %v\n", durationVerify)
//	fmt.Printf("Proof saved to %s\n", filename)
//}
