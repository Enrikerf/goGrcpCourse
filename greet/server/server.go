package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

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

func (*server) GreetManyTimes(req *proto.GreetManyTimesRequest, stream proto.GreetService_GreetManyTimesServer) error {
	fmt.Printf("greet many times function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		response := &proto.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(response)
		time.Sleep(100 * time.Millisecond)
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
	proto.RegisterGreetServiceServer(s, &server{})

	if error := s.Serve(listener); error != nil {
		log.Fatalf("fatal")
	}
}
