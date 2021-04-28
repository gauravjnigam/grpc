package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"grpc/calc/calcpb"
)

type calcServer struct{}

func (*calcServer) CalcSum(context context.Context, request *calcpb.SumRequest) (response *calcpb.SumResponse, err error) {
	log.Printf("Recieved request - %v", request)
	num1 := request.GetSumMessage().GetNum1()
	num2 := request.GetSumMessage().GetNum2()
	sum := num1 + num2
	res := calcpb.SumResponse{
		Result: sum,
	}

	return &res, nil

}

func (*calcServer) PrimeNumberDecomposition(request *calcpb.PrimeNumRequest, streamData calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Recived PrimeNumber Decomposition request %v \n", request)
	primeNum := request.GetPrimeNumMessage().GetPrimeNum()

	var k int64 = 2
	n := primeNum
	for n > 1 {
		if n%k == 0 {
			res := &calcpb.PrimeNumResponse{
				Result: k,
			}
			streamData.Send(res)
			n = n / k
			time.Sleep(1000 * time.Millisecond)
		} else {
			k = k + 1
		}
	}
	fmt.Printf("Prime number decomposition is completed by server!!!\n")
	return nil
}

func (*calcServer) ComputeAverage(streamData calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage request received !!!\n")
	result := 0.0
	sum := 0
	count := 0
	for {
		req, err := streamData.Recv()

		if err == io.EOF {
			result = float64(sum) / float64(count)
			return streamData.SendAndClose(&calcpb.AverageResponse{
				Result: result,
			})

		}
		if err != nil {
			log.Print("Error while reading the stream ", err)
			break
		}

		sum = sum + int(req.GetNum())
		count++
	}
	fmt.Printf("ComputeAverage request completed !!!\n")
	return nil
}

func (*calcServer) FindMaximum(streamData calcpb.CalculatorService_FindMaximumServer) error {
	return nil
}

func main() {
	fmt.Printf("Calculator server is starting ...\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50002")

	if err != nil {
		log.Fatalf("Server : Faild to listen - %v", err)
	}

	s := grpc.NewServer()

	calcpb.RegisterCalculatorServiceServer(s, &calcServer{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Server : failed to serve - %v", err)
	}
}
