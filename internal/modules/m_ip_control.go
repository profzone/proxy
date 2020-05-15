package modules

import "github.com/valyala/fasthttp"

type IPController struct {
	origin  []byte
	filters []ipFilter
}

func (c *IPController) scan() error {
	return nil
}

func (c *IPController) filter(req *fasthttp.Request) bool {
	if len(c.filters) > 0 {
		for _, f := range c.filters {
			if !f.filter(req) {
				return false
			}
		}
	}
	return true
}

func newIPController(rules string) (*IPController, error) {
	c := &IPController{origin: []byte(rules)}
	err := c.scan()
	if err != nil {
		return nil, err
	}

	return c, nil
}

type ipFilter interface {
	filter(req *fasthttp.Request) bool
}

type whitelistFilter struct {
	rules []string
}

func (w whitelistFilter) filter(req *fasthttp.Request) bool {
	panic("implement me")
}

type blacklistFilter struct {
	rules []string
}

func (b blacklistFilter) filter(req *fasthttp.Request) bool {
	panic("implement me")
}
