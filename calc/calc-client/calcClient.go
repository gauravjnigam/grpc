package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	"grpc/calc/calcpb"
)

func main() {
	fmt.Println("Hello Calculator service client")

	conn, err := grpc.Dial("localhost:50002", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("client is not able to connect : %v", err)
	}

	defer conn.Close()

	c := calcpb.NewCalculatorServiceClient(conn)

	//fmt.Printf("created client : %f", c)

	//result := callCalcSumService(c, 3, 10)
	//fmt.Printf("Sum from server : %d \n", result)

	//callPrimeNumDecoposition(c)

	callAverageNum(c)

	//callCurrentMax(c)

}

func callAverageNum(client calcpb.CalculatorServiceClient) {
	fmt.Printf("Calling Average number calculator\n")
	requests := []*calcpb.AverageRequest{
		&calcpb.AverageRequest{
			Num: 1,
		},
		&calcpb.AverageRequest{
			Num: 2,
		},
		&calcpb.AverageRequest{
			Num: 3,
		},
		&calcpb.AverageRequest{
			Num: 4,
		},
	}
	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatal("Error while sending computeAverage request ", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req : %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Print("Error while sending stream ", err)
	}

	log.Printf("Result : Average val = ", res.GetResult())

}

func callCurrentMax(client calcpb.CalculatorServiceClient) {
	fmt.Printf("Calling current max calculator \n")
}

func callPrimeNumDecoposition(client calcpb.CalculatorServiceClient) {
	fmt.Printf("Calling prime number decomposition \n")
	req := &calcpb.PrimeNumRequest{
		PrimeNumMessage: &calcpb.PrimeNumMessage{
			PrimeNum: 120,
		},
	}

	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition request - %v", err)
	}

	for {
		res, err := resStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming response - %v", err)
		}
		fmt.Printf("%d,", res.GetResult())

	}

	fmt.Printf("\nClient has completed the processing!!!\n")

}

func callCalcSumService(client calcpb.CalculatorServiceClient, n1 int64, n2 int64) int64 {
	req := &calcpb.SumRequest{
		SumMessage: &calcpb.SumMessage{
			Num1: n1,
			Num2: n2,
		},
	}

	res, err := client.CalcSum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator service - %v", err)
	}

	return res.GetResult()

	//return 10

}
