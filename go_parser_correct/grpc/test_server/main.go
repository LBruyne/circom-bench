package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/google/uuid"
	pb "go_parser_correct/grpc/node"
	"go_parser_correct/utils"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type server struct {
	pb.UnimplementedServiceServer
}

var (
	pk          groth16.ProvingKey
	vk          groth16.VerifyingKey
	ccs         constraint.ConstraintSystem
	NumOutput   uint32
	NumInPublic uint32
)

const wtnsParentPath = "generate_wtns/"

func generateUniqueFilename(base string, ext string, id string) string {
	return fmt.Sprintf("%s_%s.%s", base, id, ext)
}
func saveInputAsJSON(inputs []int, filename string) error {
	data := map[string][]int{
		"inputs": inputs,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	return err
}

func (s *server) Prove(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received input: %v", in.GetInput())
	if len(in.GetInput()) > 1 {
		//var inputs []int
		//for _, str := range in.GetInput() {
		//	num, err := strconv.Atoi(str)
		//	if err != nil {
		//		log.Printf("Error converting input to integer: %v", err)
		//		return &pb.Response{
		//			Code:     -1,
		//			Response: []string{"Execution failed: invalid input format"},
		//		}, nil
		//	}
		//	inputs = append(inputs, num)
		//}
		id := uuid.New().String()
		var err error
		jsonFilename := generateUniqueFilename("input", "json", id)
		err = ioutil.WriteFile("generate_wtns/input/"+jsonFilename, []byte(in.Input), 0644)
		if err != nil {
			log.Printf("Error saving input as JSON: %v", err)
			return &pb.Response{
				Code:     -1,
				Response: []string{"Execution failed: could not save input as JSON"},
			}, nil
		}
		//if err := saveInputAsJSON(inputs, "generate_wtns/input/"+jsonFilename); err != nil {
		//	log.Printf("Error saving input as JSON: %v", err)
		//	return &pb.Response{
		//		Code:     -1,
		//		Response: []string{"Execution failed: could not save input as JSON"},
		//	}, nil
		//}
		wtns_path := generateUniqueFilename("output", "wtns", id)
		//cmd := exec.CommandContext(ctx, "node", "./generate_wtns/generate_witness.js", "generate_wtns/poseidon_16_1.wasm", "generate_wtns/input.json", "generate_wtns/wtns/"+wtns_path)
		cmd := exec.Command("node", "./generate_wtns/generate_witness.js", "generate_wtns/poseidon_16_1.wasm", "generate_wtns/input/"+jsonFilename, "generate_wtns/wtns/"+wtns_path)
		if err := cmd.Run(); err != nil {
			log.Printf("Error running command: %v", err)
			return &pb.Response{
				Code:     -1,
				Response: []string{"Execution failed: could not generate witness"},
			}, nil
		}

		var w utils.R1CSCircuit

		w.Witness, w.WitnessPublic, err = utils.ParseWtns("generate_wtns/wtns/"+wtns_path, NumOutput, NumInPublic)
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
		proof, err := groth16.Prove(ccs, pk, secretWitness)
		if err != nil {
			panic(err)
		}
		err = groth16.Verify(proof, vk, witnessPublic)
		if err != nil {
			panic(err)
		}
		Result := utils.RunProveExportResult(proof, witnessPublic)
		return &pb.Response{
			Code:     0,
			Response: Result,
		}, nil
	}
	return &pb.Response{
		Code:     -1,
		Response: []string{"Execution failed due to invalid input"},
	}, nil
}

func main() {
	r1cs_path := "./generate_wtns/poseidon_16_1.r1cs"
	var err error
	ccs, NumOutput, NumInPublic, err = utils.ReadR1CS(r1cs_path)
	if err != nil {
		panic(err)
	}
	startSetup := time.Now()
	dataDir := "./key"
	pkFile, err := os.Open(dataDir + "/" + "pk.key")
	if err != nil {
		panic(err)
	}
	pk = groth16.NewProvingKey(ecc.BN254)
	bufReader := bufio.NewReaderSize(pkFile, 1024*1024)
	pk.UnsafeReadFrom(bufReader)
	defer pkFile.Close()
	vkFile, err := os.Open(dataDir + "/" + "vk.key")

	vk = groth16.NewVerifyingKey(ecc.BN254)
	vk.ReadFrom(vkFile)
	defer vkFile.Close()
	if err != nil {
		panic(err)
	}
	durationSetup := time.Since(startSetup)
	fmt.Printf("Setup time: %v\n", durationSetup)
	//——————————
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
