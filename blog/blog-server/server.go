package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"grpc/blog/blogpb"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id.omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

var collection *mongo.Collection

func (*server) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Print("Creat blog req")
	blog := request.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error : %v", err),
		)
	}
	oid := res.InsertedID.(primitive.ObjectID)
	log.Print("Creat blog req - completed!!")
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.AuthorId,
			Title:    blog.AuthorId,
			Content:  blog.Content,
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, request *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	log.Print("Reading blog : ", request.BlogId)
	blogId := request.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while parsing id "),
		)
	}
	blogData := &blogItem{}
	query := &bson.M{"_id": oid}
	result := collection.FindOne(ctx, query)
	if result.Err() != nil {
		log.Fatal("Error while reading collection from mongodb")
	}

	err = result.Decode(blogData)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Can't find blog with given ID : %v\n", err),
		)

	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blogData.ID.String(),
			AuthorId: blogData.AuthorID,
			Title:    blogData.Title,
			Content:  blogData.Content,
		},
	}, nil

}

/*
func (*server) UpdateBlog(ctx context.Context, request *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {

}
func (*server) DeleteBlog(ctx context.Context, request *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {

}
func (*server) ListBlog(ctx context.Context, request *blogpb.ListBlogRequest) (*blogpb.ListBlogResponse, error) {

}
*/

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Printf("Blog service is starting ...\n")

	// connect to mongodb
	fmt.Printf("connecting to mongodb...\n")
	client, merr := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if merr != nil {
		log.Fatalf("Server : Faild to connect to mongodb - %v", merr)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cerr := client.Connect(ctx)
	if cerr != nil {
		log.Fatal(cerr)
	}
	collection = client.Database("mydb").Collection("blog")

	//test collection
	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	if err != nil {
		log.Fatal(cerr)
	}
	id := res.InsertedID
	fmt.Printf("Inserted id : %v\n", id)

	lis, err := net.Listen("tcp", "0.0.0.0:50002")

	if err != nil {
		log.Fatalf("Server : Faild to listen - %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("Server : failed to serve - %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Stoping server...")
	s.Stop()
	fmt.Println("Closing the listner")
	lis.Close()
	fmt.Println("Stopping mongodb connection ...")
	client.Disconnect(context.TODO())

	fmt.Println("End of program!!")
}
