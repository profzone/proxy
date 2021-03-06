package models

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
