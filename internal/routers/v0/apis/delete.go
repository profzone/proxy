package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(DeleteApi{}))
}

// 删除集群
type DeleteApi struct {
	httpx.MethodDelete
}

func (req DeleteApi) Path() string {
	return ""
}

func (req DeleteApi) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
