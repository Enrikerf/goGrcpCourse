package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("hello I'm a client")

	tlsEnabled := true
	options := grpc.WithInsecure()
	if tlsEnabled {
		certFile := "../../ssl/ca.crt"
		credentials, sslError := credentials.NewClientTLSFromFile(certFile, "")
		if sslError != nil {
			log.Fatalf("error: %v", sslError)
		}
		options = grpc.WithTransportCredentials(credentials)
	}
	connection, error := grpc.Dial("localhost:50051", options)
	if error != nil {
		log.Fatalf("error: %v", error)
	}
	defer connection.Close()

	client := proto.NewGreetServiceClient(connection)
	doUnary(client)
	// doServerStreaming(client)
	// doClientStreaming(client)
	// doBiDirectional(client)
	// doUnaryWithDeadline(client, 5*time.Second)
	// doUnaryWithDeadline(client, 1*time.Second)

}

func doUnary(client proto.GreetServiceClient) {
	request := &proto.GreetRequest{
		Greeting: &proto.Greeting{
			FirstName: "name",
			LastName:  "last",
		},
	}

	response, error := client.Greet(context.Background(), request)
	if error != nil {
		log.Fatalf("error %v", error)
	}
	log.Printf("response %v", response)
}

func doServerStreaming(client proto.GreetServiceClient) {
	log.Printf("starting streaming rcp ")
	request := &proto.GreetManyTimesRequest{
		Greeting: &proto.Greeting{
			FirstName: "name",
			LastName:  "last",
		},
	}

	responseStream, error := client.GreetManyTimes(context.Background(), request)
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
		msg.GetResult()
		log.Printf("response %v", msg.GetResult())
	}
}

func doClientStreaming(client proto.GreetServiceClient) {
	fmt.Println("starting client streaming rpc")
	requests := []*proto.LongGreetRequest{
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "pepe",
			},
		},
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "john",
			},
		},
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "manuel",
			},
		},
	}
	stream, error := client.LongGreet(context.Background())
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
	fmt.Println("response %v", response)
}

func doBiDirectional(client proto.GreetServiceClient) {
	fmt.Println("bidi")

	// we create stream by invoking the client
	stream, error := client.GreetEveryone(context.Background())
	if error != nil {
		log.Fatalf("error")
		return
	}
	requests := []*proto.GreetEveryoneRequest{
		&proto.GreetEveryoneRequest{
			Greeting: &proto.Greeting{
				FirstName: "pepe",
			},
		},
		&proto.GreetEveryoneRequest{
			Greeting: &proto.Greeting{
				FirstName: "john",
			},
		},
		&proto.GreetEveryoneRequest{
			Greeting: &proto.Greeting{
				FirstName: "manuel",
			},
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
			fmt.Println("received: %v", response.GetResult())
		}
		close(waitChannel)
	}()
	//block until everyhting is done
	<-waitChannel
}

func doUnaryWithDeadline(client proto.GreetServiceClient, timeout time.Duration) {
	request := &proto.GreetWithDeadLineRequest{
		Greeting: &proto.Greeting{
			FirstName: "name",
			LastName:  "last",
		},
	}
	context, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, error := client.GreetWithDeadLine(context, request)
	if error != nil {
		statusError, ok := status.FromError(error)
		if ok {
			if statusError.Code() == codes.DeadlineExceeded {
				fmt.Println("timeout %v", error)
			}
		} else {
			log.Fatalf("unexpected %v", error)
		}
		return
	}
	log.Printf("response %v", response)
}
