package main

import (
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/frontend"
	"go_parser_correct/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	main1()
	//main2()
	main3()
}

func main1() {

	baseDir := "../test_poseidon"
	subDir := "300"
	jsFile := "poseidon_16_1_js/generate_witness.js"
	wasmFile := "poseidon_16_1_js/poseidon_16_1.wasm"
	inputFile := "../test_poseidon/input.json"
	outputFile := "./utils/output.wtns"

	// 构建文件路径
	jsFilePath := fmt.Sprintf("%s/%s/%s", baseDir, subDir, jsFile)
	wasmFilePath := fmt.Sprintf("%s/%s/%s", baseDir, subDir, wasmFile)

	// 构建命令
	cmd := exec.Command("node", jsFilePath, wasmFilePath, inputFile, outputFile)
	//cmd := exec.Command("node", "../test_poseidon/generate_witness.js", "../test_poseidon/poseidon_16_1.wasm", "../test_poseidon/input.json", "./utils/output.wtns")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", output)
	//filePath := "./utils/output.wtns"
	//witnesses, err := utils.ParseWtns(filePath)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//
	//for i, witness := range witnesses {
	//	fmt.Printf("Witness %d: %v\n", i, witness)
	//}
}

//func generateWitness(a int) {
//
//	cmd := exec.Command("node", "../test_poseidon/generate_witness.js", "../test_poseidon/poseidon_16_1.wasm", "../test_poseidon/input.json", "./utils/output.wtns")
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//		return
//	}
//	fmt.Printf("Output: %s\n", output)
//}

// 安装参数保存读取的测试版本
func main2() {

	runtime.GOMAXPROCS(4)
	dir, _ := os.Getwd()
	startRead := time.Now()
	fmt.Println("working directory: ", dir)
	ccs, err := ReadR1CS("r1cs")
	durationRead := time.Since(startRead)
	fmt.Printf("Read time: %v\n", durationRead)
	if err != nil {
		panic(err)
	}
	a, b, c := ccs.GetNbVariables()
	fmt.Println(a, b, c)
	var w R1CSCircuit
	startParse := time.Now()
	//w.Witness, err = utils.ParseWtns("./output.wtns")
	w.Witness, err = utils.ParseWtns("./utils/output.wtns")
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
	//pk, vk, err := groth16.Setup(ccs)
	_, _, err = groth16.Setup(ccs)
	if err != nil {
		panic(err)
	}
	durationSetup := time.Since(startSetup)
	fmt.Printf("Setup time: %v\n", durationSetup)
	//存储的Pk
	//data, err := json.Marshal(pk)
	//if err != nil {
	//	fmt.Println("Error serializing:", err)
	//	return
	//}
	//err = ioutil.WriteFile("Pk.json", data, 0644)
	//if err != nil {
	//	fmt.Println("Error writing to file:", err)
	//	return
	//}
	ReadJsonTime := time.Now()
	dataFromFile, err := ioutil.ReadFile("Pk.json")
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}
	var deserializedPK *groth16_bn254.ProvingKey
	err = json.Unmarshal(dataFromFile, &deserializedPK)
	if err != nil {
		fmt.Println("Error deserializing:", err)
		return
	}
	durationReadJson := time.Since((ReadJsonTime))
	fmt.Printf("ReadPKJson time: %v\n", durationReadJson)
	// 存储的Vk
	//dataVk, err := json.Marshal(vk)
	//if err != nil {
	//	fmt.Println("Error serializing:", err)
	//	return
	//}
	//err = ioutil.WriteFile("Vk.json", dataVk, 0644)
	//if err != nil {
	//	fmt.Println("Error writing to file:", err)
	//	return
	//}
	ReadJsonTime = time.Now()
	dataFromFile, err = ioutil.ReadFile("Vk.json")
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}
	var deserializedVK *groth16_bn254.VerifyingKey
	err = json.Unmarshal(dataFromFile, &deserializedVK)
	if err != nil {
		fmt.Println("Error deserializing:", err)
		return
	}
	durationReadJson = time.Since((ReadJsonTime))
	fmt.Printf("ReadVKJson time: %v\n", durationReadJson)

	//↑↑↑↑↑↑↑↑↑————————————————数据加载——————————————————↑↑↑↑↑↑↑↑↑

	startProve := time.Now()
	//proof, err := groth16.Prove(ccs, pk, secretWitness)
	proof, err := groth16.Prove(ccs, deserializedPK, secretWitness, backend.WithIcicleAcceleration())
	if err != nil {
		panic(err)
	}
	durationProve := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", durationProve)
	startVerify := time.Now()
	_ = groth16.Verify(proof, deserializedVK, witnessPublic)
	durationVerify := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", durationVerify)

}

// 原始版本的测试
func main3() {

	runtime.GOMAXPROCS(4)
	dir, _ := os.Getwd()
	startRead := time.Now()
	fmt.Println("working directory: ", dir)
	ccs, err := ReadR1CS("r1cs_300")
	durationRead := time.Since(startRead)
	fmt.Printf("Read time: %v\n", durationRead)
	if err != nil {
		panic(err)
	}
	a, b, c := ccs.GetNbVariables()
	fmt.Println(a, b, c)
	var w R1CSCircuit
	w.Witness, err = utils.ParseWtns("./utils/output.wtns")
	if err != nil {
		panic(err)
	}
	//w.Witness = make([]frontend.Variable, 11093)
	//for i := 0; i < len(w.Witness); i++ {
	//	w.Witness[i] = frontend.Variable(0)
	//}

	//scs := ccs.(*cs.SparseR1CS)
	//srs, srsLagrange, err := unsafekzg.NewSRS(scs)
	//if err != nil {
	//	panic(err)
	//}
	//witnessFull, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
	//if err != nil {
	//	panic(err)
	//}
	//witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
	//if err != nil {
	//	panic(err)
	//}
	//pk, vk, err := plonk.Setup(ccs, srs, srsLagrange)
	////_, err := plonk.Setup(r1cs, kate, &publicWitness)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//proof, err := plonk.Prove(ccs, pk, witnessFull)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = plonk.Verify(proof, vk, witnessPublic)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//
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

	startProve := time.Now()
	//proof, err := groth16.Prove(ccs, pk, secretWitness)
	proof, err := groth16.Prove(ccs, pk, secretWitness, backend.WithIcicleAcceleration())
	if err != nil {
		panic(err)
	}
	durationProve := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", durationProve)
	startVerify := time.Now()
	_ = groth16.Verify(proof, vk, witnessPublic)
	durationVerify := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", durationVerify)
}
