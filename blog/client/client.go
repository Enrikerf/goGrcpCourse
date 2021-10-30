package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"proto"
)

func main() {
	fmt.Println("hello I'm a Blog client")

	options := grpc.WithInsecure()
	connection, err := grpc.Dial("localhost:50051", options)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer connection.Close()

	client := proto.NewBlogServiceClient(connection)
	blog := proto.Blog{
		AuthorId: "Enrique",
		Title:    "Linea de costa",
		Content:  "1,2 Apisonadora",
	}
	createdBlog, err := client.CreateBlog(context.Background(), &proto.CreateBlogRequest{Blog: &blog})
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("blog %v", createdBlog)
}
