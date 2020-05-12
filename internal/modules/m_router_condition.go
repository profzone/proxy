package modules

import "github.com/valyala/fasthttp"

type routerCondition struct {
	rules []rule
}

func newRouterCondition(condition string) *routerCondition {
	if condition == "" {
		return nil
	}
	return &routerCondition{}
}

func (c *routerCondition) Match(req *fasthttp.Request) bool {
	for _, r := range c.rules {
		if !r.match(req) {
			return false
		}
	}
	return true
}

func (c *routerCondition) UnmarshalJSON([]byte) error {
	panic("implement me")
}

func (c routerCondition) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (c *routerCondition) GobDecode([]byte) error {
	panic("implement me")
}

func (c routerCondition) GobEncode() ([]byte, error) {
	panic("implement me")
}
