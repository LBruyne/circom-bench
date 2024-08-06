package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	fr_bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"io/ioutil"
	"math/big"
)

//func generateWitnessAndParse(wasmFile, inputFile, outputWitnessFile string) (*R1CSCircuit, error) {
//	// 执行Node.js脚本生成witness
//	cmd := exec.Command("node", "generate_witness.js", wasmFile, inputFile, outputWitnessFile)
//	err := cmd.Run()
//	if err != nil {
//		return nil, fmt.Errorf("failed to run Node.js script: %v", err)
//	}
//
//	// 读取生成的witness文件
//	witnessData, err := ioutil.ReadFile(outputWitnessFile)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read witness file: %v", err)
//	}
//
//	// 解析witness数据到Go结构体
//	var circuit R1CSCircuit
//	circuit.Witness, err = parseWitness(witnessData)
//	if err != nil {
//		return nil, fmt.Errorf("failed to parse witness data: %v", err)
//	}
//
//	return &circuit, nil
//}

func ParseWtns(filePath string, NumOutput uint32, NumInPublic uint32) ([]frontend.Variable, []frontend.Variable, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	if string(fileContent[:4]) != "wtns" {
		return nil, nil, fmt.Errorf("invalid file format")
	}

	//version := binary.LittleEndian.Uint32(fileContent[4:8])
	//fmt.Println("Version:", version)
	//
	//sections := binary.LittleEndian.Uint32(fileContent[8:12])
	//fmt.Println("Number of sections:", sections)
	//
	//idSection1 := binary.LittleEndian.Uint32(fileContent[12:16])
	//fmt.Println("ID Section 1:", idSection1)
	//
	//idSection1Length := binary.LittleEndian.Uint32(fileContent[16:20])
	//fmt.Println("ID Section 1 Length (first 32 bits):", idSection1Length)
	//
	//idSection1LengthHigh := binary.LittleEndian.Uint32(fileContent[20:24])
	//fmt.Println("ID Section 1 Length (last 32 bits):", idSection1LengthHigh)
	//
	//n32 := binary.LittleEndian.Uint32(fileContent[24:28])
	//fmt.Println("n32:", n32)

	rawPrimeStart := 28
	rawPrimeEnd := rawPrimeStart + int(8)*4
	//rawPrime := fileContent[rawPrimeStart:rawPrimeEnd]
	//fmt.Println("Raw Prime:", rawPrime)

	witnessSize := binary.LittleEndian.Uint32(fileContent[rawPrimeEnd : rawPrimeEnd+4])
	fmt.Println("Witness Size:", witnessSize)

	witnesses := make([]frontend.Variable, 0, witnessSize-NumInPublic-NumOutput)
	witnessesPublic := make([]frontend.Variable, NumInPublic+NumOutput)

	idSection2Start := rawPrimeEnd + 4
	//idSection2 := binary.LittleEndian.Uint32(fileContent[idSection2Start : idSection2Start+4])
	//fmt.Println("ID Section 2:", idSection2)
	//
	//idSection2Length := binary.LittleEndian.Uint32(fileContent[idSection2Start+4 : idSection2Start+8])
	//fmt.Println("ID Section 2 Length (first 32 bits):", idSection2Length)
	//
	//idSection2LengthHigh := binary.LittleEndian.Uint32(fileContent[idSection2Start+8 : idSection2Start+12])
	//fmt.Println("ID Section 2 Length (last 32 bits):", idSection2LengthHigh)

	witnessDataStart := idSection2Start + 12
	witnessDataLength := int(8 * 4)
	for i := uint32(0); i < witnessSize; i++ {
		witnessStart := witnessDataStart + int(i)*witnessDataLength
		witnessEnd := witnessStart + witnessDataLength
		witness := fileContent[witnessStart:witnessEnd]
		//witnesses[i] = new(big.Int).SetBytes(witness)
		witness = convertLittleEndianToBigEndian(witness)
		bigIntWitness := new(big.Int).SetBytes(witness)

		if i >= 1 && i < NumOutput+1+NumInPublic {
			witnessesPublic[i-1] = bigIntWitness
			fmt.Printf("%d\n", bigIntWitness)
			//↑用于观测具体的witness是否输入outpu和publicinput
		} else {
			witnesses = append(witnesses, bigIntWitness)
		}
	}

	return witnesses, witnessesPublic, nil
}

func convertLittleEndianToBigEndian(data []byte) []byte {
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-1-i] = data[len(data)-1-i], data[i]
	}
	return data
}
func RunProveExportResult(proof groth16.Proof, witnessPublic witness.Witness) []string {
	_proof, ok := proof.(interface{ MarshalSolidity() []byte })
	if !ok {
		panic("proof does not implement MarshalSolidity()")
	}
	proofStr := hex.EncodeToString(_proof.MarshalSolidity())
	bPublicWitness, err := witnessPublic.MarshalBinary()
	if err != nil {
		panic(err)
	}
	bPublicWitness = bPublicWitness[12:]
	publicWitnessStr := hex.EncodeToString(bPublicWitness)
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
	var __proof [8]*big.Int
	for i := 0; i < 8; i++ {
		__proof[i] = new(big.Int).SetBytes(proofBytes[fpSize*i : fpSize*(i+1)])
	}
	var finalproof []string
	var finalinput []string
	for _, bi := range __proof {
		finalproof = append(finalproof, bi.String())
	}
	for _, bi := range input {
		finalinput = append(finalinput, bi.String())
	}
	fmt.Println(finalproof, "||", finalinput)
	var result []string
	result = append(result, finalproof...)
	result = append(result, finalinput...)
	return result
}

//func main() {
//	filePath := "./utils/output.wtns"
//	witnesses, err := ParseWtns(filePath)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	for i, witness := range witnesses {
//		fmt.Printf("Witness %d: %v\n", i, witness)
//	}
//}
