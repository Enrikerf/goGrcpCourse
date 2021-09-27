package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a client")

	connection, error := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if error != nil {
		log.Fatalf("error: %v", error)
	}
	defer connection.Close()

	client := proto.NewGreetServiceClient(connection)
	doUnary(client)
	doServerStreaming(client)
	doClientStreaming(client)

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
	fmt.Printf("response %v", response)
}
