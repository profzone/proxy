package modules

import (
	"bytes"
	"encoding/gob"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/pkg/route"
)

type Router struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
	// 路由条件
	Condition string `json:"condition" default:""`
	// URL重写规则
	RewritePattern string `json:"rewritePattern" default:""`
	// 重写到特定集群
	ClusterID  uint64 `json:"clusterID,string" default:""`
	conditions *routerCondition
}

func (r *Router) GobDecode(data []byte) error {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	err := dec.Decode(r)
	if err != nil {
		return err
	}

	r.conditions = newRouterCondition(r.Condition)
	return nil
}

func (r *Router) Match(req *fasthttp.Request, params route.Params) bool {
	if r.conditions != nil && !r.conditions.Match(req, params) {
		return false
	}
	return true
}

func (r *Router) Rewrite(req *fasthttp.Request, params route.Params) error {
	if r.RewritePattern == "" {
		return nil
	}
	expr := newRewriteExpr(req, r.RewritePattern, params)
	if expr.Error() != nil {
		return expr.Error()
	}
	err := expr.apply()
	if err != nil {
		return err
	}

	req.SetRequestURI(expr.uri())
	return nil
}
