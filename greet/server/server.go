package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"proto"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, in *proto.GreetRequest) (*proto.GreetResponse, error) {
	fmt.Printf("greet function was invoked with %v\n", in)
	firstName := in.GetGreeting().GetFirstName()
	result := "hello" + firstName
	respose := &proto.GreetResponse{
		Result: result,
	}
	return respose, nil
}

func main() {
	fmt.Println("Hello world")

	listener, error := net.Listen("tcp", "0.0.0.0:50051")
	if error != nil {
		log.Fatalf("failed to listen")
	}
	s := grpc.NewServer()
	proto.RegisterGreetServiceServer(s, &server{})

	if error := s.Serve(listener); error != nil {
		log.Fatalf("fatal")
	}
}
