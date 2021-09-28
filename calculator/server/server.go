package main

import (
	"calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Printf("average invoked")
	var nNumbers int32 = 0
	var sum int32 = 0
	for {
		request, error := stream.Recv()
		if error == io.EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Average: sum / nNumbers,
			})
		}
		if error != nil {
			log.Fatalf("error")
		}
		sum += request.GetNumber()
		nNumbers++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Finding maximum")
	var currentMaximum int32 = 0
	for {
		request, error := stream.Recv()
		if error == io.EOF {
			return nil
		}
		if error != nil {
			log.Fatalf("error")
			return error
		}

		if request.GetNumber() > currentMaximum {
			currentMaximum = request.GetNumber()
			sendError := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: currentMaximum,
			})
			if sendError != nil {
				log.Fatalf("error")
				return sendError
			}
		}
	}
}

func (*server) SquareRoot(context context.Context, request *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("received squareroot rpc")
	number := request.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("received negative number: %v", number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
