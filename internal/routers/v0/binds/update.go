package binds

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(UpdateBind{}))
}

// 更新绑定信息
type UpdateBind struct {
	httpx.MethodPatch
	// 编号
	ID   uint64          `name:"id,string" in:"path"`
	Body modules.Cluster `name:"body" in:"body"`
}

func (req UpdateBind) Path() string {
	return "/:id"
}

func (req UpdateBind) Output(ctx context.Context) (result interface{}, err error) {
	req.Body.ID = req.ID
	err = modules.UpdateCluster(&req.Body, storage.Database)
	return
}
