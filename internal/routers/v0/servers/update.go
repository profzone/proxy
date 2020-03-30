package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(UpdateServer{}))
}

// 更新集群
type UpdateServer struct {
	httpx.MethodPatch
}

func (req UpdateServer) Path() string {
	return ""
}

func (req UpdateServer) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
