package modules

import (
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
	ClusterID uint64 `json:"clusterID,string" default:""`
}

func (r *Router) Match(req *fasthttp.Request, params route.Params) bool {
	condition := newRouterCondition(r.Condition)
	if condition != nil && !condition.Match(req) {
		return false
	}
	return true
}

func (r *Router) Rewrite(req *fasthttp.Request, params route.Params) error {
	if r.RewritePattern == "" {
		return nil
	}
	expr := newRewriteExpr(req, r.RewritePattern, nil)
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
