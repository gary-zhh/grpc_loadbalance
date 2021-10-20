package main

import (
	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"

	"context"
	"hello_grpc/sum_grpc"
	"log"
	"net"
	"time"
)
type server struct {
	sum_grpc.UnimplementedSumServer
}

func (s *server) Add(ctx context.Context, in *sum_grpc.AddRequest) (*sum_grpc.AddReply, error) {
	return &sum_grpc.AddReply{
		Sum: in.A + in.B,
	}, nil
}

func main(){
	conn, err := net.DialTimeout("tcp", ":8089", time.Second*5)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}
	log.Println("TCP connection success. local addr: ",conn.LocalAddr())
	if _, err = conn.Write([]byte("sum")); err != nil {
		log.Fatal("transport algorithm name error, please check")
	}

	srvConn, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		log.Fatalf("couldn't create yamux server: %s", err)
	}

	grpcServer := grpc.NewServer()
	sum_grpc.RegisterSumServer(grpcServer,&server{})
	log.Println("listening to grpc client over TCP connection...")
	if err := grpcServer.Serve(srvConn); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
