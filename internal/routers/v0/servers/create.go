package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(CreateServer{}))
}

// 创建集群
type CreateServer struct {
	httpx.MethodPost
}

func (req CreateServer) Path() string {
	return ""
}

func (req CreateServer) Output(ctx context.Context) (result interface{}, err error) {
	panic("implement me")
}
