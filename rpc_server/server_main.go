package main

import (
	"context"
	pb "hello_grpc/hello_grpc"
	"hello_grpc/teststream"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + in.Name}, nil
}

type streamServer struct {
	teststream.UnimplementedTestStreamServer
}

func (s *streamServer) Send(r teststream.TestStream_SendServer) error {
	for {
		req, err := r.Recv()
		if err != nil {
			return err
		}
		for i := range req.Req {
			req.Req[i]++
		}
		err = r.Send(&teststream.Response{
			Res: req.Req,
		})
		if err != nil {
			return err
		}
	}
}
func main() {
	// 要监听的协议和端口
	lis, err := net.Listen("tcp", ":9898")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 实例化gRPC server结构体
	s := grpc.NewServer()
	//服务注册
	pb.RegisterGreeterServer(s, &server{})
	log.Println("开始监听，等待远程调用...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//teststream.RegisterTestStreamServer(s, &streamServer{})
	//log.Println("开始监听，等待远程调用...")
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}
