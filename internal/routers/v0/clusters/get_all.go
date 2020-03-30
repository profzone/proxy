package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(GetClusters{}))
}

// 获取所有集群
type GetClusters struct {
	httpx.MethodGet
}

func (req GetClusters) Path() string {
	return ""
}

func (req GetClusters) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
