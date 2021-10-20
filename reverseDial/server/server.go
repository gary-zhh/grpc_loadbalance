package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	pb "hello_grpc/hello_grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + in.Name}, nil
}

type InstanceInfo struct {

}

type AlgoServer interface {
	//Invoke(ctx context.Context,req interface{}) (interface{}, error)
	//HeartBeat(ctx context.Context) InstanceInfo
}

type Listener struct {
	rawConn net.Conn
}
var _ net.Listener= (*Listener)(nil)
func (l *Listener)Accept()(net.Conn,error) {
	if l.rawConn == nil {
		return nil, errors.New("no rawConn")
	}
	return l.rawConn, nil
}
func (l *Listener)Close()error {
	return l.rawConn.Close()
}
func (l *Listener)Addr()net.Addr {
	return l.rawConn.LocalAddr()
}
type ActiveGrpcServer struct {
	//*grpc.ServiceDesc
	*grpc.Server
	endpoint string
	listener net.Listener
}

func New(desc *grpc.ServiceDesc, algo AlgoServer, endpoint string) {
	server := &ActiveGrpcServer{
		Server: grpc.NewServer(),
		endpoint: endpoint,
	}
	server.RegisterService(desc, algo)
	if conn, err := net.Dial("tcp", endpoint); err != nil {
		panic("panic")
	} else {
		fmt.Println("send connection success", conn.LocalAddr())
		server.listener = &Listener{
			rawConn: conn,
		}
	}
	//go func() {
		err := server.Serve(server.listener)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	//}()
}


func main() {
	ch := make(chan struct{})
	go func(){
		l, err := net.Listen("tcp",":9898")
		if err != nil {
			log.Fatal(err)
		}
		close(ch)
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receive connection success", conn.LocalAddr())
	}()
	<-ch
	New(&pb.Greeter_ServiceDesc,&server{},":9898")
}
