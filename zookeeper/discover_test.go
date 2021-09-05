package zookeeper

import (
	"context"
	"fmt"
	pb "hello_grpc/hello_grpc"
	_ "hello_grpc/loadbalance/weight"
	"testing"

	"google.golang.org/grpc/resolver"

	"google.golang.org/grpc"
)

func Test_Discovery(t *testing.T) {
	r := NewSvcDiscovery([]string{"127.0.0.1:2181"})
	resolver.Register(r)
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///server", r.Scheme()),
		grpc.WithInsecure(),
		grpc.WithBalancerName("weight"))
	fmt.Println(err)
	client := pb.NewGreeterClient(conn)
	rep, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
	rep, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
	rep, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
	rep, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
	conn.Close()
}
