package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(CreateCluster{}))
}

// 创建集群
type CreateCluster struct {
	httpx.MethodPost
}

func (req CreateCluster) Path() string {
	return ""
}

func (req CreateCluster) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
