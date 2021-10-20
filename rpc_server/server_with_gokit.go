package main

import (
	"context"
	pb "hello_grpc/hello_grpc"
	"log"
	"net"

	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
)

type serverWithGoKti struct {
	pb.UnimplementedGreeterServer
	endpoint.Endpoint
}

func Reply(req string) (res string) {
	return "hello, " + req
}

type strFunc func(string) string

func makeEndpoint(f strFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res := make(chan interface{})
		go func() {
			req := request.(*pb.HelloRequest)
			res <- f(req.GetName())
		}()
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		case ret := <-res:
			return &pb.HelloReply{Message: ret.(string)}, nil
		}
	}
}

func (s *serverWithGoKti) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	resp, err := s.Endpoint(ctx, in)
	return resp.(*pb.HelloReply), err
}

func main() {
	// 要监听的协议和端口
	lis, err := net.Listen("tcp", ":9898")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	helloServer := &serverWithGoKti{Endpoint: makeEndpoint(Reply)}
	// 实例化gRPC server结构体
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, helloServer)
	log.Println("开始监听，等待远程调用...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
