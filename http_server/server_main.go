package main

import (
	"context"
	"flag"
	"fmt"
	pb "hello_grpc/hello_grpc"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	echoEndpoint = flag.String("echo_endpoint", ":9898", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":9999", mux)
}
func main() {
	if err := run(); err != nil {
		fmt.Print(err.Error())
	}
}
