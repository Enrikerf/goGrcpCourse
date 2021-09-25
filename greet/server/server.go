package main

import (
	"fmt"
	"log"
	"net"

	"proto"

	"google.golang.org/grpc"
)

type server struct{}

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
