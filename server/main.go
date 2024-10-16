package main

import (
	"context"
	"fmt"
	"grpc-lt/pb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTmpServiceServer
}

func (*server) UnaryRPC(ctx context.Context, req *pb.TmpRequest) (*pb.TmpResponse, error) {
	fmt.Println("--------------------")
	fmt.Println(string(req.GetMsg()))

	msg := "Hello Client!!"
	res := &pb.TmpResponse{
		Msg: msg,
	}

	return res, nil
}

func (*server) ServerStreamingRPC(req *pb.TmpRequest, stream pb.TmpService_ServerStreamingRPCServer) error {
	fmt.Println("--------------------")
	fmt.Println(string(req.GetMsg()))

	for i := 1; i <= 100; i++ {
		res := &pb.TmpResponse{Msg: "Hello Client" + strconv.Itoa(i)}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

func (*server) ClientStreamingRPC(stream pb.TmpService_ClientStreamingRPCServer) error {
	fmt.Println("--------------------")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.TmpResponse{Msg: "Hello Client!!"}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		fmt.Println(req.GetMsg())
	}
}

func (*server) BidirectionalStreamingRPC(stream pb.TmpService_BidirectionalStreamingRPCServer) error {
	fmt.Println("--------------------")

	cnt := 0
	for {
		cnt++
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println(req.GetMsg())

		err = stream.Send(&pb.TmpResponse{Msg: "Hello Client" + strconv.Itoa(cnt)})
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTmpServiceServer(s, &server{})

	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve:%v", err)
	}
}
