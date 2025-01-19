package wrr

import (
	"context"
	"math"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const Name = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{
		HealthCheck: true,
	})
}

func init() {
	balancer.Register(newBuilder())
}

type PickerBuilder struct {
	picker *Picker
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if p.picker == nil {
		conns := make([]*weightConn, 0, len(info.ReadySCs))
		for sc, sci := range info.ReadySCs {
			md, _ := sci.Address.Metadata.(map[string]interface{})
			weightVal := md["weight"]
			weight, _ := weightVal.(float64)
			conns = append(conns, &weightConn{
				SubConn:       sc,
				weight:        int(weight),
				currentWeight: int(weight),
			})
		}
		p.picker = &Picker{
			conns: conns,
		}
	}
	return p.picker
}

type Picker struct {
	conns []*weightConn
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var (
		totalWeight int
		res         *weightConn
	)
	for _, c := range p.conns {
		c.mutex.Lock()
		totalWeight += c.efficientWeight
		c.currentWeight += c.efficientWeight
		if res == nil || res.currentWeight < c.currentWeight {
			res = c
		}
		c.mutex.Unlock()
	}
	res.mutex.Lock()
	res.currentWeight -= totalWeight
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.SubConn,
		Done: func(di balancer.DoneInfo) {
			res.mutex.Lock()
			defer res.mutex.Unlock()
			if di.Err != nil && res.efficientWeight == 0 {
				return
			}
			switch di.Err {
			case nil:
				if res.efficientWeight == math.MaxUint32 {
					return
				}
				// 增加权重
				res.efficientWeight++
			case context.DeadlineExceeded:
				// 超时可以考虑动态调整
				// 比如，第一次超时降低 1，第二次超时降低 2，第三次超时降低 4
				// 第 n 次超时降低 2^(n-1)
				res.efficientWeight -= 10
			default:
				code := status.Code(di.Err)
				switch code {
				case codes.Unauthenticated:
					// 熔断
					res.efficientWeight = 1
				case codes.ResourceExhausted:
					res.efficientWeight >>= 1
				case codes.Aborted:
					// 降级
					res.efficientWeight >>= 1
				default:
					if res.efficientWeight == 1 {
						// 降无可降
						return
					}
					res.efficientWeight--
				}
			}
		},
	}, nil
}

type weightConn struct {
	balancer.SubConn
	mutex           sync.Mutex
	weight          int
	currentWeight   int
	efficientWeight int
}
