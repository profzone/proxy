package modules

import (
	"github.com/valyala/fasthttp"
	"sync/atomic"
)

type LoadBalancer interface {
	Apply(req fasthttp.Request, servers []ServerContract) uint64
}

// RoundRobin round robin loadBalance impl
type RoundRobin struct {
	ops *uint64
}

// NewRoundRobin create a RoundRobin
func NewRoundRobin() LoadBalancer {
	var ops uint64
	ops = 0

	return RoundRobin{
		ops: &ops,
	}
}

// Select select a server from servers using RoundRobin
func (rr RoundRobin) Apply(req fasthttp.Request, servers []ServerContract) uint64 {
	l := uint64(len(servers))

	if 0 >= l {
		return 0
	}

	target := servers[int(atomic.AddUint64(rr.ops, 1)%l)]
	return target.GetIdentity()
}
