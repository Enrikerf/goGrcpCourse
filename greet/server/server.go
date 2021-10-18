package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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

func (*server) LongGreet(stream proto.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet invoked")
	result := "Hello"
	for {
		request, error := stream.Recv()
		if error == io.EOF {
			return stream.SendAndClose(&proto.LongGreetResponse{
				Result: result,
			})
		}
		if error != nil {
			log.Fatalf("error")
		}
		firstName := request.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream proto.GreetService_GreetEveryoneServer) error {
	fmt.Println("Hi everyone")

	for {
		request, error := stream.Recv()
		if error == io.EOF {
			return nil
		}
		if error != nil {
			log.Fatalf("error")
			return error
		}
		result := "Hello " + request.GetGreeting().GetFirstName() + "! "
		sendError := stream.Send(&proto.GreetEveryoneResponse{
			Result: result,
		})
		if sendError != nil {
			log.Fatalf("error")
			return sendError
		}
	}

}

func (*server) GreetWithDeadLine(ctx context.Context, request *proto.GreetWithDeadLineRequest) (*proto.GreetWithDeadLineResponse, error) {
	fmt.Printf("greet with deadline function was invoked with %v\n", request)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("client canceled the request")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := request.GetGreeting().GetFirstName()
	result := "hello nodeadline: " + firstName
	respose := &proto.GreetWithDeadLineResponse{
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
	certificate := "../../ssl/server.crt"
	keyFile := "../../ssl/server.pem"
	credentials, sslError := credentials.NewServerTLSFromFile(certificate, keyFile)
	if sslError != nil {
		log.Fatalf("error: %v", sslError)
		return
	}
	options := grpc.Creds(credentials)
	s := grpc.NewServer(options)
	proto.RegisterGreetServiceServer(s, &server{})

	if error := s.Serve(listener); error != nil {
		log.Fatalf("fatal")
	}
}
