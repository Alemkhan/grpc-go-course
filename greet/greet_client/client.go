package main

import (
	"com.grpc.tleu/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	calcService := greetpb.NewCalculatorServiceClient(conn)
	runPrimeNumberDecomposition(120, calcService)
	runComputeAverage([]int32{1, 2, 3, 4}, calcService)
}

func runPrimeNumberDecomposition(n int32, calcService greetpb.CalculatorServiceClient) {
	req := &greetpb.NumberRequest{Number: n}
	stream, err := calcService.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatal("Error 500")
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("Error with response from server stream RPC %v", err)
		}
		log.Printf(fmt.Sprint(res.GetResult(), " "))
	}
}

func runComputeAverage(numbersStream []int32, calcService greetpb.CalculatorServiceClient) {
	stream, err := calcService.ComputerAverage(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to server")
	}
	for _, n := range numbersStream {
		stream.Send(&greetpb.NumberRequest{Number: n})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error recieving response")
	}

	fmt.Printf("Result: %f", res.GetResult())
}
