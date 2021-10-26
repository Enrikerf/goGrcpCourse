package main

import (
	"fmt"
	"log"
	"net"

	"github.com/enrikerf/grcpGoProof/blog/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Hello world")

	listener, error := net.Listen("tcp", "0.0.0.0:50051")
	if error != nil {
		log.Fatalf("failed to listen")
	}
	tlsEnabled := true
	var serverOptions grpc.ServerOption
	if tlsEnabled {
		certificate := "../../ssl/server.crt"
		keyFile := "../../ssl/server.pem"
		credentials, sslError := credentials.NewServerTLSFromFile(certificate, keyFile)
		if sslError != nil {
			log.Fatalf("error: %v", sslError)
			return
		}
		serverOptions = grpc.Creds(credentials)
	}

	s := grpc.NewServer(serverOptions)
	proto.Register(s, &server{})

	if error := s.Serve(listener); error != nil {
		log.Fatalf("fatal")
	}
}
