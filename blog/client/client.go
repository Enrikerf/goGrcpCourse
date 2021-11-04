package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
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
	//insertBlog(client)
	//readBlog(client)
	//updateBlog(client)
	//deleteBlog(client)
	listBlogs(client)
}

func listBlogs(client proto.BlogServiceClient) {
	request := &proto.ListBlogRequest{}
	responseStream, err := client.ListBlog(context.Background(), request)

	if err != nil {
		log.Fatalf("error %v", err)
	}
	for {
		msg, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("response %v", err)
		}
		log.Printf("response %v", msg)
	}
}

func deleteBlog(client proto.BlogServiceClient) {
	deleteBlogRequest := proto.DeleteBlogRequest{
		BlogId: "617d03b51677ccfc7ec5a9bd",
	}
	deleteBlogResponse, err := client.DeleteBlog(context.Background(), &deleteBlogRequest)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("blog %v", deleteBlogResponse)
}

func updateBlog(client proto.BlogServiceClient) {
	blog := proto.Blog{
		Id: "617d03b51677ccfc7ec5a9bd",
		AuthorId: "Enrique",
		Title:    "Linea de costa",
		Content:  "no apisona",
	}
	createdBlog, err := client.UpdateBlog(context.Background(), &proto.UpdateBlogRequest{Blog: &blog})
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("blog %v", createdBlog)
}

func insertBlog(client proto.BlogServiceClient) {

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

func readBlog(client proto.BlogServiceClient)  {
	blogRequest := proto.ReadBlogRequest{
		BlogId: "617d03b51677ccfc7ec5a9bd",
	}
	readBlog, err := client.ReadBlog(context.Background(), &blogRequest)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("blog %v", readBlog)

}
