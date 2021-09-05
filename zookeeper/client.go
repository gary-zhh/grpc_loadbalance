package zookeeper

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type zkCli struct {
	conn      zk.Conn
	endpoints sync.Map
}

func callback(ev zk.Event) {
	fmt.Println()
	fmt.Println("#############################")
	fmt.Println("type = ", ev.Type)
	fmt.Println("path = ", ev.Path)
	fmt.Println("state = ", ev.State)
	fmt.Println("server = ", ev.Server)
	fmt.Println("#############################")
	fmt.Println()
}
func New(endpoints []string, duration time.Duration) *zkCli {
	//eventCallbackOption := zk.WithEventCallback(callback)
	conn, _, err := zk.Connect(endpoints, duration)
	if err != nil {
		return nil
	}
	return &zkCli{
		conn: *conn,
	}
}

func (cli *zkCli) Close() {
	cli.conn.Close()
}
func (cli *zkCli) Get(path string) string {
	b, _, err := cli.conn.Get(path)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
func (cli *zkCli) GetChildren(path string) []string {
	chs, _, err := cli.conn.Children(path)
	if err != nil {
		return nil
	}
	for i := range chs {
		chs[i] = strings.TrimSuffix(path, "/") + "/" + chs[i]
	}
	return chs
}
func (cli *zkCli) WatchChildren(ctx context.Context, path string) <-chan zk.Event {
	events := make(chan zk.Event)
	go func() {
		defer close(events)
		for {
			chs, _, evs, _ := cli.conn.ChildrenW(path)
			for _, c := range chs {
				childpath := strings.TrimSuffix(path, "/") + "/" + c
				if _, ok := cli.endpoints.Load(childpath); !ok {
					cli.endpoints.Store(childpath, struct{}{})
					fmt.Println("##### store ", childpath)
					go cli.WatchNode(ctx, childpath, events)
					select {
					case <-ctx.Done():
						return
					case events <- zk.Event{
						Type:  zk.EventNodeCreated,
						Path:  childpath,
						State: zk.StateUnknown,
					}:
					}
				}

			}
			select {
			case <-ctx.Done():
				return
			case <-evs:
			}
		}
	}()

	return events
}

func (cli *zkCli) WatchNode(ctx context.Context, path string, events chan<- zk.Event) {
	fmt.Println("##### watch node ", path)
	for {
		_, _, evs, _ := cli.conn.ExistsW(path)
		select {
		case <-ctx.Done():
			return
		case ev := <-evs:
			events <- ev
			if ev.Type == zk.EventNodeDeleted {
				fmt.Println("##### delete node ", ev.Path)
				cli.endpoints.Delete(ev.Path)
				return
			}
		}
	}
}
