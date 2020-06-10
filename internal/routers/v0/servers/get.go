package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(GetServer{}))
}

// 获取单个服务器
type GetServer struct {
	httpx.MethodGet

	ID uint64 `name:"id" in:"path"`
}

func (req GetServer) Path() string {
	return "/:id"
}

func (req GetServer) Output(ctx context.Context) (result interface{}, err error) {
	server, err := modules.GetServer(req.ID, storage.Database)
	if err != nil {
		return
	}

	result = server
	return
}
