package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(UpdateApi{}))
}

// 更新集群
type UpdateApi struct {
	httpx.MethodPatch
}

func (req UpdateApi) Path() string {
	return ""
}

func (req UpdateApi) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
