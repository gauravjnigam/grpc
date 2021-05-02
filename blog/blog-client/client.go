package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"grpc/blog/blogpb"
)

func main() {
	fmt.Println("Hello Calculator service client")

	conn, err := grpc.Dial("localhost:50002", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("client is not able to connect : %v\n", err)
	}

	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	// create blog client
	blogId := createBlogRequest(c)

	// read blog clinet
	readBlogRequest(c, blogId)

}

func readBlogRequest(client blogpb.BlogServiceClient, blogId string) {
	_, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "",
	})

	if err != nil {
		fmt.Printf("Error happened during reading the request - %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: blogId,
	}
	res, err := client.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Printf("Error happened during reading the request - %v\n", err)
	}

	fmt.Printf("Read blog : %v\n", res)

}

func createBlogRequest(client blogpb.BlogServiceClient) (blogId string) {

	blog := &blogpb.Blog{
		AuthorId: "Gaurav N",
		Title:    "gn blog",
		Content:  "gn blog content",
	}

	res, err := client.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error while creating blog - %v\n", err)
	}
	fmt.Printf("Blog has been created - %v\n", res.Blog.AuthorId)

	return res.Blog.Id
}
