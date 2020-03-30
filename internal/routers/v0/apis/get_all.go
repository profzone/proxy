package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(GetApis{}))
}

// 获取所有集群
type GetApis struct {
	httpx.MethodGet
}

func (req GetApis) Path() string {
	return ""
}

func (req GetApis) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
