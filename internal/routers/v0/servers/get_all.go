package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(GetServers{}))
}

// 获取所有集群
type GetServers struct {
	httpx.MethodGet
}

func (req GetServers) Path() string {
	return ""
}

func (req GetServers) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
