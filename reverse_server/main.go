package main

import (
	"context"
	"fmt"
	pb "hello_grpc/hello_grpc"
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

type Listener struct {
	conn net.Conn
}

func (l *Listener) Accept() (net.Conn, error) {
	return l.conn, nil
}

func (l *Listener) Close() error {
	return nil
}

func (l *Listener) Addr() net.Addr {
	return l.conn.LocalAddr()
}
func get() net.Listener {
	conn, err := net.Dial("tcp", ":9898")
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return &Listener{
		conn: conn,
	}
}

func main() {
	go func() {
		// 要监听的协议和端口
		lis, err := net.Listen("tcp", ":9898")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("failed to accept: %v", err)
		}
		fmt.Println("local addr:", conn.LocalAddr())
	}()

	// 实例化gRPC server结构体
	s := grpc.NewServer()
	// 服务注册
	pb.RegisterGreeterServer(s, &server{})
	log.Println("开始监听，等待远程调用...")
	if err := s.Serve(get()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
