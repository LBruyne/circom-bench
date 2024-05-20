package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/consensys/gnark/frontend"
	"io/ioutil"
)

type R1CSCircuit struct {
	Witness []frontend.Variable
}

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

func ParseWtns(filePath string) ([]frontend.Variable, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if string(fileContent[:4]) != "wtns" {
		return nil, fmt.Errorf("invalid file format")
	}

	version := binary.LittleEndian.Uint32(fileContent[4:8])
	fmt.Println("Version:", version)

	sections := binary.LittleEndian.Uint32(fileContent[8:12])
	fmt.Println("Number of sections:", sections)

	idSection1 := binary.LittleEndian.Uint32(fileContent[12:16])
	fmt.Println("ID Section 1:", idSection1)

	idSection1Length := binary.LittleEndian.Uint32(fileContent[16:20])
	fmt.Println("ID Section 1 Length (first 32 bits):", idSection1Length)

	idSection1LengthHigh := binary.LittleEndian.Uint32(fileContent[20:24])
	fmt.Println("ID Section 1 Length (last 32 bits):", idSection1LengthHigh)

	n32 := binary.LittleEndian.Uint32(fileContent[24:28])
	fmt.Println("n32:", n32)

	rawPrimeStart := 28
	rawPrimeEnd := rawPrimeStart + int(8)*4
	rawPrime := fileContent[rawPrimeStart:rawPrimeEnd]
	fmt.Println("Raw Prime:", rawPrime)

	witnessSize := binary.LittleEndian.Uint32(fileContent[rawPrimeEnd : rawPrimeEnd+4])
	fmt.Println("Witness Size:", witnessSize)

	witnesses := make([]frontend.Variable, witnessSize)

	idSection2Start := rawPrimeEnd + 4
	idSection2 := binary.LittleEndian.Uint32(fileContent[idSection2Start : idSection2Start+4])
	fmt.Println("ID Section 2:", idSection2)

	idSection2Length := binary.LittleEndian.Uint32(fileContent[idSection2Start+4 : idSection2Start+8])
	fmt.Println("ID Section 2 Length (first 32 bits):", idSection2Length)

	idSection2LengthHigh := binary.LittleEndian.Uint32(fileContent[idSection2Start+8 : idSection2Start+12])
	fmt.Println("ID Section 2 Length (last 32 bits):", idSection2LengthHigh)

	witnessDataStart := idSection2Start + 12
	witnessDataLength := int(8 * 4)
	for i := uint32(0); i < witnessSize; i++ {
		witnessStart := witnessDataStart + int(i)*witnessDataLength
		witnessEnd := witnessStart + witnessDataLength
		witness := fileContent[witnessStart:witnessEnd]
		//witnesses[i] = new(big.Int).SetBytes(witness)
		witnesses[i] = witness
	}

	return witnesses, nil
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
