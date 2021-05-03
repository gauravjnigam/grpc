package main

import (
	"context"
	"fmt"
	"io"
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

	// read blog client
	readBlogRequest(c, blogId)

	// update blog client
	//updateBlogRequest(c, blogId)

	// delete the blog client
	//deleteBlogRequest(c, blogId)

	// list the blog client
	listBlogRequest(c)
}

func listBlogRequest(client blogpb.BlogServiceClient) {
	log.Print("List Blog request...")
	stream, err := client.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatal("Error while streaming the ListBlog response")
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error while streaming the response ")
		}

		fmt.Printf("Received blog - %v\n", res.GetBlog())
	}

}

func deleteBlogRequest(client blogpb.BlogServiceClient, blogId string) {
	log.Printf("Deleting the blog - %s", blogId)

	res, err := client.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogId,
	})

	if err != nil {
		fmt.Printf("Error while updating the blog id : %v", err)
	}

	fmt.Printf("Blog - %s is deleted!! %v", blogId, res)

}
func updateBlogRequest(client blogpb.BlogServiceClient, blogId string) {
	log.Printf("Updating the blog - %s", blogId)
	updatedblog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Ravi N",
		Title:    "rn blog",
		Content:  "rn blog content",
	}

	res, err := client.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: updatedblog,
	})

	if err != nil {
		log.Fatal("Error while updating the blog")
	}
	fmt.Printf("Blog was updated - %v\n", res)
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
