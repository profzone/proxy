package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(UpdateApi{}))
}

// 更新API
type UpdateApi struct {
	httpx.MethodPatch
	// 编号
	ID   uint64      `name:"id,string" in:"path"`
	Body modules.API `name:"body" in:"body"`
}

func (req UpdateApi) Path() string {
	return "/:id"
}

func (req UpdateApi) Output(ctx context.Context) (result interface{}, err error) {
	req.Body.ID = req.ID
	err = modules.UpdateAPI(&req.Body, storage.Database)
	return
}
