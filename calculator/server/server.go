package main

import (
	"calculatorpb"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("sum function was invoked with %v\n", in)

	respose := &calculatorpb.SumResponse{
		SumResult: in.GetFirstNumber() + in.GetSecondNumber(),
	}
	return respose, nil
}

func (*server) PrimeDecomposition(request *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	fmt.Printf("Prime Decomposition was invoked with %v\n", request)
	primeNumber := request.GetPrimeNumber()
	var k int32 = 2
	for primeNumber > 1 {
		if primeNumber%k == 0 {
			response := &calculatorpb.PrimeDecompositionResponse{
				PrimeFactor: primeNumber,
			}
			stream.Send(response)
			primeNumber = primeNumber / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello world")

	listener, error := net.Listen("tcp", "0.0.0.0:50051")
	if error != nil {
		log.Fatalf("failed to listen")
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if error := s.Serve(listener); error != nil {
		log.Fatalf("fatal")
	}
}
