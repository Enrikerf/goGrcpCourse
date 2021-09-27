package main

import (
	"calculatorpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a client")

	connection, error := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if error != nil {
		log.Fatalf("error: %v", error)
	}
	defer connection.Close()

	client := calculatorpb.NewCalculatorServiceClient(connection)
	doUnary(client)
	doPrimeDecomposition(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	request := &calculatorpb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 2,
	}

	response, error := client.Sum(context.Background(), request)
	if error != nil {
		log.Fatalf("error %v", error)
	}
	log.Printf("response %v", response)
}

func doPrimeDecomposition(client calculatorpb.CalculatorServiceClient) {
	request := &calculatorpb.PrimeDecompositionRequest{
		PrimeNumber: 94840284,
	}

	responseStream, error := client.PrimeDecomposition(context.Background(), request)

	if error != nil {
		log.Fatalf("error %v", error)
	}
	for {
		msg, error := responseStream.Recv()
		if error == io.EOF {
			break
		}
		if error != nil {
			log.Printf("response %v", error)
		}
		log.Printf("response %v", msg)
	}
}