package modules

import (
	"github.com/valyala/fasthttp"
	"testing"
)

func TestScanPattern(t *testing.T) {
	req := fasthttp.AcquireRequest()
	defer func() {
		fasthttp.ReleaseRequest(req)
	}()
	req.URI().QueryArgs().Add("key", "123")
	req.URI().QueryArgs().Add("version", "v1")

	pattern1 := "/peer/v0/rewrite?key=$(origin.query.key)&version=$(origin.query.version)&target=abc"
	real1 := "/peer/v0/rewrite?key=123&version=v1&target=abc"
	expr1 := newRewriteExpr(req, pattern1, nil)
	expr1.apply()

	if real1 != expr1.uri() {
		t.Errorf("%s is not equals real1 %s", expr1.uri(), real1)
	}

	pattern2 := "/peer/v0/rewrite?key=$(origin.query.key)$(origin.query.version)&target=abc"
	real2 := "/peer/v0/rewrite?key=123v1&target=abc"
	expr2 := newRewriteExpr(req, pattern2, nil)
	expr2.apply()

	if real2 != expr2.uri() {
		t.Errorf("%s is not equals real2 %s", expr2.uri(), real2)
	}
}
