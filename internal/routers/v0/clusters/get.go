package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(GetCluster{}))
}

// 获取单个集群
type GetCluster struct {
	httpx.MethodGet

	ID uint64 `name:"id" in:"path"`
}

func (req GetCluster) Path() string {
	return "/:id"
}

func (req GetCluster) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
