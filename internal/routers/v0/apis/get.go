package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(GetApi{}))
}

// 获取单个集群
type GetApi struct {
	httpx.MethodGet

	ID uint64 `name:"id" in:"path"`
}

func (req GetApi) Path() string {
	return "/:id"
}

func (req GetApi) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
