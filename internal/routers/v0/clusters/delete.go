package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(DeleteCluster{}))
}

// 删除集群
type DeleteCluster struct {
	httpx.MethodDelete
}

func (req DeleteCluster) Path() string {
	return ""
}

func (req DeleteCluster) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
