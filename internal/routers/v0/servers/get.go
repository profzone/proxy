package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(GetServer{}))
}

// 获取单个集群
type GetServer struct {
	httpx.MethodGet

	ID uint64 `name:"id" in:"path"`
}

func (req GetServer) Path() string {
	return "/:id"
}

func (req GetServer) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
