package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(CreateApi{}))
}

// 创建集群
type CreateApi struct {
	httpx.MethodPost
}

func (req CreateApi) Path() string {
	return ""
}

func (req CreateApi) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
