package modules

import (
	"github.com/valyala/fasthttp"
	"testing"
)

func TestScanRule(t *testing.T) {
	req := fasthttp.AcquireRequest()
	req.URI().QueryArgs().Add("version", "2")
	req.URI().QueryArgs().Add("count", "5")
	req.URI().QueryArgs().Add("versionStr", "b")
	req.URI().QueryArgs().Add("countStr", "d")

	condition1 := "query.version>=1,query.count<10"
	cond1 := newRouterCondition(condition1, nil)
	if !cond1.Match(req) {
		t.Errorf("!cond1.Match(req)")
	}

	condition2 := "query.versionStr@>=a,query.countStr@<f"
	cond2 := newRouterCondition(condition2, nil)
	if !cond2.Match(req) {
		t.Errorf("!cond2.Match(req)")
	}
}
