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

	rep, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
	rep, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhou"})
	fmt.Println(rep, err)
	//client := teststream.NewTestStreamClient(conn)
	//i := byte(0)
	//for {
	//	req := &teststream.Request{Req: []byte{i}}
	//	i++
	//	s, err := client.Send(context.Background())
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = s.Send(req)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	res, err := s.Recv()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println(res.Res)
	//	time.Sleep(time.Second * 5)
	//}

}
