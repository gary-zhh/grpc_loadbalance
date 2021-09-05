package weight

import (
	"math/rand"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

const Name = "weight"

func init() {
	balancer.Register(newBuilder())
}
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: false})
}

type attKey struct {
}

type Weight int

func SetWeight(addr resolver.Address, weight Weight) resolver.Address {
	addr.Attributes = addr.Attributes.WithValues(attKey{}, weight)
	return addr
}

func GetWeight(addr resolver.Address) Weight {
	v := addr.Attributes.Value(attKey{})
	ret, _ := v.(Weight)
	return ret
}

type PickerBuilder struct {
}

type Picker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
}

const (
	minWeight = 1
	maxWeight = 5
)

func (*PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("weightPicker: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for subConn, addr := range info.ReadySCs {
		node := GetWeight(addr.Address)
		if node <= 0 {
			node = minWeight
		} else if node > 5 {
			node = maxWeight
		}
		for i := 0; i < int(node); i++ {
			scs = append(scs, subConn)
		}
	}
	return &Picker{
		subConns: scs,
	}
}

func (p *Picker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	index := rand.Intn(len(p.subConns))
	sc := p.subConns[index]
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}
