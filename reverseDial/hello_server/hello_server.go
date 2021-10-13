package main

import (
	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"
	"hello_grpc/hello_grpc"

	"context"
	"log"
	"net"
	"time"
)
type server struct {
	hello_grpc.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *hello_grpc.HelloRequest) (*hello_grpc.HelloReply, error) {
	return &hello_grpc.HelloReply{
		Message: "hello, " + in.Name,
	}, nil
}

func main(){
	conn, err := net.DialTimeout("tcp", ":8089", time.Second*5)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}
	log.Println("TCP connection success. local addr: ",conn.LocalAddr())
	if _, err = conn.Write([]byte("hello")); err != nil {
		log.Fatal("transport algorithm name error, please check")
	}

	srvConn, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		log.Fatalf("couldn't create yamux server: %s", err)
	}

	grpcServer := grpc.NewServer()
	hello_grpc.RegisterGreeterServer(grpcServer,&server{})
	log.Println("listening to grpc client over TCP connection...")
	if err := grpcServer.Serve(srvConn); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
