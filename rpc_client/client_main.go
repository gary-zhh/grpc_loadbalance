package main

import (
	"context"
	"fmt"
	pb "hello_grpc/hello_grpc"
	_ "hello_grpc/loadbalance/weight"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":9898", grpc.WithInsecure(), grpc.WithBalancerName("weight"))
	fmt.Println(err)
	client := pb.NewGreeterClient(conn)
	rep, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
}
