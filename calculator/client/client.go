package main

import (
	"calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	doClientStreaming(client)
	doFindMaximumBidirectional(client)
	doErrorUnary(client)
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

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting client streaming rpc")
	requests := []*calculatorpb.AverageRequest{
		&calculatorpb.AverageRequest{
			Number: 1,
		},
		&calculatorpb.AverageRequest{
			Number: 2,
		},
		&calculatorpb.AverageRequest{
			Number: 3,
		},
	}
	stream, error := client.Average(context.Background())
	if error != nil {
		log.Fatalf("error")
	}
	for _, request := range requests {
		fmt.Printf("sending: %v\n", request)
		stream.Send(request)
	}

	response, error := stream.CloseAndRecv()
	if error != nil {
		log.Fatalf("error")
	}
	fmt.Printf("response: %v", response)
}

func doFindMaximumBidirectional(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("find maximum")

	// we create stream by invoking the client
	stream, error := client.FindMaximum(context.Background())
	if error != nil {
		log.Fatalf("error")
		return
	}
	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 3,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 4,
		},
	}
	waitChannel := make(chan struct{})
	// we send a bunch of messages to de client go routin
	go func() {
		for _, request := range requests {
			fmt.Println("sending message %v", request)
			stream.Send(request)
			time.Sleep(1000 * time.Microsecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages form the client go routine
	go func() {
		for {
			response, error := stream.Recv()
			if error == io.EOF {
				break
			}
			if error != nil {
				log.Fatalf("error")
				break
			}
			fmt.Println("received: %v", response)
		}
		close(waitChannel)
	}()
	//block until everyhting is done
	<-waitChannel
}

func doErrorUnary(client calculatorpb.CalculatorServiceClient) {
	request := &calculatorpb.SquareRootRequest{
		Number: -10,
	}

	response, error := client.SquareRoot(context.Background(), request)
	if error != nil {
		responseError, ok := status.FromError(error)
		if ok {
			fmt.Println(responseError.Message())
			fmt.Println(responseError.Code())
			if responseError.Code() == codes.InvalidArgument {
				fmt.Println("we sent negative number")
			}
		} else {
			log.Fatalln("error %v", error)
		}

	}
	log.Printf("response %v", response)
}
