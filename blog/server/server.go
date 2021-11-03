package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"proto"
)

//make it global
var collection *mongo.Collection

type Server struct {
}

func (server *Server) ReadBlog(ctx context.Context, request *proto.ReadBlogRequest) (*proto.ReadBlogResponse, error) {
	fmt.Println("Read request")
	blogIdString := request.GetBlogId()
	objectId, err := primitive.ObjectIDFromHex(blogIdString)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("cannot parse id"))
	}
	data := &blogItem{}
	filter := bson.D{{"_id",objectId}}
	result := collection.FindOne(context.Background(),filter)
	if err := result.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("not found"))
	}

	return &proto.ReadBlogResponse{Blog: &proto.Blog{
		Id: data.ID.String(),
		AuthorId: data.AuthorID,
		Content: data.Content,
		Title: data.Title,
	}}, err


}

func (server *Server) CreateBlog(ctx context.Context, request *proto.CreateBlogRequest) (*proto.CreateBlogResponse, error) {
	fmt.Println("create blog request received")
	blog := request.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	result, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error $v", err))
	}

	objectId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("can't parse $v", ok))
	}
	return &proto.CreateBlogResponse{
		Blog: &proto.Blog{
			Id:       objectId.Hex(),
			AuthorId: blog.GetTitle(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil

}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {
	// if we crash the go code, we ge the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Hello world")

	// Mongo setup ---
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Blog Service Started")
	collection = client.Database("go_grcp_course_db").Collection("blog")
	//ctx:= context.TODO()
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	//err = client.Ping(ctx, readpref.Primary())
	//if err = client.Disconnect(ctx); err != nil {
	//	fmt.Println("panic")
	//	panic(err)
	//}
	//
	//collection := client.Database("go_grcp_course").Collection("blog")
	//fmt.Println("%v",collection)
	// Mongo setup --------

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen")
	}

	var serverOptions []grpc.ServerOption
	server := grpc.NewServer(serverOptions...)
	proto.RegisterBlogServiceServer(server, &Server{})

	go func() {
		fmt.Println("Starting Server...")
		if err := server.Serve(listener); err != nil {
			log.Fatalf("fatal")
		}
	}()

	// Wait for control C to exit
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Bock until a signal is received
	<-channel
	fmt.Println("Stopping the server")
	server.Stop()
	fmt.Println("closing the listener")
	listener.Close()
	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")
}
