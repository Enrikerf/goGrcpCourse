package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"proto"

	"google.golang.org/grpc"
)

type Server struct {
}

func main() {
	fmt.Println("Hello world")

	listener, error := net.Listen("tcp", "0.0.0.0:50051")
	if error != nil {
		log.Fatalf("failed to listen")
	}

	serverOptions := []grpc.ServerOption{}
	server := grpc.NewServer(serverOptions...)
	proto.RegisterBlogServiceServer(server, &Server{})

	go func() {
		fmt.Println("Starting Server...")
		if error := server.Serve(listener); error != nil {
			log.Fatalf("fatal")
		}
	}()

	// Wait for control C to exit
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Bock until a signal is received
	<-channel
	server.Stop()
	fmt.Println("closing the listener")
	listener.Close()
}
