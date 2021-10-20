package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/yamux"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"hello_grpc/hello_grpc"
	"hello_grpc/sum_grpc"
	"log"
	"net"
	"sync"
)
type algorithm struct {
	name string
	id string
	stub interface{}
}

type Client struct {
	servers map[string][]algorithm
	listener net.Listener
	done chan struct{}
	mu sync.Mutex
}
func (c *Client) Run() {
	for {
		select {
		case <-c.done:
			return
		default:
			log.Println("waiting for incoming TCP connections...")
			conn, err := c.listener.Accept()
			if err != nil {
				log.Fatalf("couldn't accept %s", err)
				return
			}
			log.Println("==================================")
			log.Println("receiving connection from ", conn.RemoteAddr())
			buf := make([]byte,1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalf("couldn't read %s", err)
			}
			algoName := string(buf[:n])
			log.Println("remote server can invoke \"", algoName, "\"")
			wrapConn, err := yamux.Client(conn, yamux.DefaultConfig())
			if err != nil {
				log.Fatalf("couldn't create yamux %s", err)
			}
			log.Println("trying to dial grpc server over incoming TCP connection")
			grpcConn, err := grpc.Dial("", grpc.WithInsecure(),
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return wrapConn.Open()
				}))
			if err != nil {
				log.Fatalf("did not connect: %s", err)
			}
			log.Println("dial grpc server success")
			log.Println("==================================")
			c.newGRPCClient(algoName,grpcConn)
		}
	}
}
func (c *Client) Close() {
	close(c.done)
	c.listener.Close()
}
func New(port int) (*Client,error){
	log.Println("launch tcp server and listening...")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d",port))
	if err != nil {
		log.Fatalf("could not listen: %s", err)
		return nil,err
	}
	client := &Client{
		servers: make(map[string][]algorithm),
		listener:ln,
		done: make(chan struct{}),
	}
	go client.Run()
	return client, nil
}
func main() {
	client,_ := New(8089)
	for {
		index := 0
		fmt.Println("[ INPUT THE ALGORITHM YOU WANT TO CALL   1.SUM     2.HELLO]")
		_,_ = fmt.Scanln(&index)
		switch index {
		case 1:
			if algos,ok := client.servers["sum"]; !ok || len(algos) == 0 {
				fmt.Println("[ WARNING!! ALGO SUM NOT PREPARED YET]")
				continue
			}
			fmt.Println("[ INPUT 2 NUMS TO GET THEIR SUM ]")
			a,b := int64(0),int64(0)
			fmt.Scanln(&a)
			fmt.Scanln(&b)
			if selected, ok := client.servers["sum"][0].stub.(sum_grpc.SumClient); !ok {
				log.Fatal("internal error")
			} else {
				res,err := selected.Add(context.Background(), &sum_grpc.AddRequest{A:a,B:b})
				if err != nil {
					log.Fatalf("error when calling Add: %s", err)
				}
				fmt.Println("[ THE SUM IS ]: ", res.Sum)
			}
		case 2:
			if algos,ok := client.servers["hello"]; !ok || len(algos) == 0 {
				fmt.Println("[ WARNING!! ALGO HELLO NOT PREPARED YET]")
				continue
			}
			fmt.Println("[ INPUT WHAT YOU WANT TO SEND ]")
			str := ""
			fmt.Scanln(&str)
			if selected, ok := client.servers["hello"][0].stub.(hello_grpc.GreeterClient); !ok {
				log.Fatal("internal error")
			} else {
				res,err := selected.SayHello(context.Background(), &hello_grpc.HelloRequest{Name: str})
				if err != nil {
					log.Fatalf("error when calling Add: %s", err)
				}
				fmt.Println("[ RESPONSE ]: ", res.Message)
			}
		}
	}
}

func (c *Client)newGRPCClient(algoName string, conn *grpc.ClientConn) {
	algo := algorithm{
		name:algoName,
		id: uuid.NewV4().String(),
	}
	switch algoName {
	case "hello":
		algo.stub = hello_grpc.NewGreeterClient(conn)
	case "sum":
		algo.stub = sum_grpc.NewSumClient(conn)
	default:
		log.Fatal("algorithm not supported")
	}
	c.mu.Lock()
	if _, ok := c.servers[algoName]; !ok {
		c.servers[algoName] = make([]algorithm,0)
	}
	c.servers[algoName] = append(c.servers[algoName], algo)
	c.mu.Unlock()
}
func handleConn() {

	//response, err := c.Add(context.Background(), &sum_grpc.AddRequest{A:1,B:2})
	//if err != nil {
	//	log.Fatalf("error when calling SayHello: %s", err)
	//}
	//log.Printf("response from server: %s", response.Sum)
}

