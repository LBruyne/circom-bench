package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	pb "go_parser_correct/grpc/node"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fileName := "example.json"
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	r, err := c.Prove(ctx, &pb.Request{Input: string(content)})
	if err != nil {
		log.Fatalf("could not execute prove: %v", err)
	}
	log.Printf("Response code: %d, Response: %v", r.GetCode(), r.GetResponse())
}
