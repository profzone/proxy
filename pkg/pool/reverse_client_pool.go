package pool

import (
	"github.com/valyala/fasthttp"
	"sync"
)

var ClientPool *ReverseClientPool

func init() {
	ClientPool = &ReverseClientPool{}
}

type ReverseClientPool struct {
	sync.Pool
}

func (p *ReverseClientPool) AcquireClient() *fasthttp.Client {
	value := ClientPool.Get()
	if value == nil {
		return &fasthttp.Client{}
	}
	return value.(*fasthttp.Client)
}

func (p *ReverseClientPool) ReleaseClient(c *fasthttp.Client) {
	if c != nil {
		ClientPool.Put(c)
	}
}
