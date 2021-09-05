package zookeeper

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_Client_All(t *testing.T) {
	cli := New([]string{"8.136.223.137:2181"}, time.Second*5)
	if cli == nil {
		t.Error("cli is nil")
	}

	chs := cli.GetChildren("/gRPC/server")
	for _, c := range chs {
		fmt.Println(c, ": ", cli.Get(c))
	}

	evs := cli.WatchChildren(context.Background(), "/test")
	// evs = cli.WatchNode(context.Background(), "/test")
	for ev := range evs {
		fmt.Println()
		fmt.Println("#############################")
		fmt.Println("type = ", ev.Type)
		fmt.Println("path = ", ev.Path)
		fmt.Println("state = ", ev.State)
		fmt.Println("server = ", ev.Server)
		fmt.Println("#############################")
		fmt.Println()
	}
}
