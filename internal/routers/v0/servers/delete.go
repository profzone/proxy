package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(DeleteServer{}))
}

// 删除集群
type DeleteServer struct {
	httpx.MethodDelete
}

func (req DeleteServer) Path() string {
	return ""
}

func (req DeleteServer) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
