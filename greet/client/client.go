package main

import (
	"fmt"

	"log"

	"proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a client")

	connection, error := grcp.Dial("localhost:50051", grpc.WithInsecure())
	if error != nil {
		log.Fatalf("error: %v", error)
	}
	defer connection.close()

	client := proto.NewGreetServiceClient(connection)

}
