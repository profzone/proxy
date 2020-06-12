package servers

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(DeleteServer{}))
}

// 删除服务器
type DeleteServer struct {
	httpx.MethodDelete
	// 编号
	ID uint64 `name:"id,string" in:"path"`
}

func (req DeleteServer) Path() string {
	return "/:id"
}

func (req DeleteServer) Output(ctx context.Context) (result interface{}, err error) {
	err = modules.DeleteServer(req.ID, storage.Database)
	return
}
