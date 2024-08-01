package main

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_encoder(t *testing.T) {
	bigNumber := new(big.Int)
	bigNumber.SetString("949386769277681657017765594388685360182081120054278437478610802394466877440", 10)
	//bigNumber.SetString("16", 10)

	// Number of bits in each uint32
	bitsPerUint32 := uint(32)

	// Slice to store the uint32 chunks
	var uint32Chunks []uint32

	// Mask to extract the lower 32 bits
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bitsPerUint32), big.NewInt(1))

	// Calculate the uint32 representation
	for i := 0; i < 8; i++ {
		// Extract the lower 32 bits
		chunk := new(big.Int).And(bigNumber, mask)
		uint32Chunks = append(uint32Chunks, uint32(chunk.Uint64()))

		// Shift the number right by 32 bits
		bigNumber.Rsh(bigNumber, bitsPerUint32)
	}

	// Since we are extracting from the least significant part first,
	// the chunks are in little-endian order and need not be reversed.
	fmt.Println("uint32 Chunks:", uint32Chunks)
	uint8Chunks := make([]uint8, 0, 32) // 长度为 32 的 uint8 数组
	for _, chunk := range uint32Chunks {
		uint8Chunks = append(uint8Chunks,
			uint8(chunk&0xFF),       // 最低有效字节
			uint8((chunk>>8)&0xFF),  // 第二个字节
			uint8((chunk>>16)&0xFF), // 第三个字节
			uint8((chunk>>24)&0xFF)) // 最高有效字节
	}
	// Step 2: Decode
	reconstructedNumber := decodeUint8ToBigInt(uint8Chunks)
	//reconstructedNumber := big.NewInt(0)
	//for i, byteVal := range uint8Chunks {
	//	tempBigInt := new(big.Int).SetUint64(uint64(byteVal))
	//	tempBigInt.Lsh(tempBigInt, uint(8*i))
	//	reconstructedNumber.Add(reconstructedNumber, tempBigInt)
	//}
	fmt.Println("The reconstructed big integer is:", reconstructedNumber)
}
func decodeUint8ToBigInt(uint8Chunks []uint8) *big.Int {
	reconstructedNumber := big.NewInt(0)
	for i, byteVal := range uint8Chunks {
		tempBigInt := new(big.Int).SetUint64(uint64(byteVal))
		tempBigInt.Lsh(tempBigInt, uint(8*i))
		reconstructedNumber.Add(reconstructedNumber, tempBigInt)
	}
	return reconstructedNumber
}
