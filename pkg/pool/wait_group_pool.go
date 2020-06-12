package pool

import (
	"sync"
)

var WGPool *WaitGroupPool

func init() {
	WGPool = &WaitGroupPool{}
}

type WaitGroupPool struct {
	sync.Pool
}

func (p *WaitGroupPool) AcquireWG() *sync.WaitGroup {
	value := WGPool.Get()
	if value == nil {
		return &sync.WaitGroup{}
	}
	return value.(*sync.WaitGroup)
}

func (p *WaitGroupPool) ReleaseWG(c *sync.WaitGroup) {
	if c != nil {
		WGPool.Put(c)
	}
}
