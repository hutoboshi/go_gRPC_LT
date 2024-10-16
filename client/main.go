package main

import (
	"context"
	"fmt"
	"grpc-lt/pb"
	"io"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect:%v", err)
	}
	defer conn.Close()

	client := pb.NewTmpServiceClient(conn)

	var input int
	for {
		fmt.Println("-----------------------------")
		fmt.Println("select gRPC")
		fmt.Println("1.Unary RPC")
		fmt.Println("2.Server Streaming RPC")
		fmt.Println("3.Client Streaming RPC")
		fmt.Println("4.Bidirectional Streaming RPC")
		fmt.Println("-----------------------------")
		fmt.Print(":")
		fmt.Scan(&input)
		fmt.Println()
		if input == 1 {
			callUnaryRPC(client)
		} else if input == 2 {
			callServerStreamingRPC(client)
		} else if input == 3 {
			callClientStreamingRPC(client)
		} else if input == 4 {
			callBidirectionalStreamingRPC(client)
		} else if input == 100 {
			fmt.Println("End")
			break
		} else {
			fmt.Println("No gRPC")
		}
	}
}

func callUnaryRPC(client pb.TmpServiceClient) {
	req := &pb.TmpRequest{Msg: "Hello Server!!"}
	res, err := client.UnaryRPC(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.GetMsg())
}

func callServerStreamingRPC(client pb.TmpServiceClient) {
	stream, err := client.ServerStreamingRPC(context.Background(), &pb.TmpRequest{Msg: "Hello Server!!"})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(res.GetMsg()))
	}
}

func callClientStreamingRPC(client pb.TmpServiceClient) {
	stream, err := client.ClientStreamingRPC(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for i := 1; i <= 100; i++ {
		sendErr := stream.Send(&pb.TmpRequest{Msg: "Hello Server" + strconv.Itoa(i)})
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(50 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetMsg())
}

func callBidirectionalStreamingRPC(client pb.TmpServiceClient) {
	stream, err := client.BidirectionalStreamingRPC(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// request
	go func() {
		for i := 1; i <= 100; i++ {
			sendErr := stream.Send(&pb.TmpRequest{Msg: "Hello Server" + strconv.Itoa(i)})
			if sendErr != nil {
				log.Fatalln(sendErr)
			}
			time.Sleep(50 * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(res.GetMsg())
		}
		close(ch)
	}()
	<-ch
}
