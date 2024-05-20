package main

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"go_parser_correct/utils"
	"os"
	"runtime"
	"time"
)

func main() {
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
	w.Witness, err = utils.ParseWtns("./output.wtns")
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
	//
	//pk, vk, err := plonk.Setup(ccs, srs, srsLagrange)
	////_, err := plonk.Setup(r1cs, kate, &publicWitness)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
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

	//似乎原版是plonkcommitments因此groth16不行

}
