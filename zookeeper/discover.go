package zookeeper

import (
	"context"
	"fmt"
	"hello_grpc/loadbalance/weight"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"google.golang.org/grpc/resolver"
)

type SvcDiscovery struct {
	zk       *zkCli
	zkEvents <-chan zk.Event
	cc       resolver.ClientConn
	services sync.Map
	prefix   string
}

var _ resolver.Builder = (*SvcDiscovery)(nil)

var _ resolver.Resolver = (*SvcDiscovery)(nil)

const (
	Schema = "gRPC"
)

func NewSvcDiscovery(endpoints []string) resolver.Builder {
	zk := New(endpoints, time.Second*5)
	if zk == nil {
		log.Fatal("nil zookeeper client")
		return nil
	}
	return &SvcDiscovery{
		zk: zk,
	}
}

func (s *SvcDiscovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	s.cc = cc
	s.prefix = fmt.Sprintf("/%s/%s", target.Scheme, target.Endpoint)
	s.ResolveNow(resolver.ResolveNowOptions{})
	go s.watcher()
	return s, nil
}
func (*SvcDiscovery) Scheme() string {
	return Schema
}
func (s *SvcDiscovery) watcher() {
	s.zkEvents = s.zk.WatchChildren(context.Background(), s.prefix)
	for ev := range s.zkEvents {
		switch ev.Type {
		case zk.EventNodeCreated, zk.EventNodeDataChanged:
			w, _ := strconv.Atoi(s.zk.Get(ev.Path))
			s.services.Store(ev.Path,
				weight.SetWeight(resolver.Address{
					Addr: strings.TrimPrefix(ev.Path, s.prefix),
				}, weight.Weight(w)))
		case zk.EventNodeDeleted:
			s.services.Delete(ev.Path)
		}
		_ = s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
	}
}

func (s *SvcDiscovery) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	s.services.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}

func (s *SvcDiscovery) ResolveNow(resolver.ResolveNowOptions) {
	children := s.zk.GetChildren(s.prefix)
	for _, c := range children {
		w, _ := strconv.Atoi(s.zk.Get(c))
		s.services.Store(c,
			weight.SetWeight(resolver.Address{
				Addr: strings.TrimPrefix(c, s.prefix+"/"),
			}, weight.Weight(w)))
	}
	_ = s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
}
func (s *SvcDiscovery) Close() {
	s.zk.Close()
}
