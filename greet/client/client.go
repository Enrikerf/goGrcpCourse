package main

import (
	"context"
	"fmt"
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
