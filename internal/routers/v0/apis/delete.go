package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(DeleteApi{}))
}

// 删除API
type DeleteApi struct {
	httpx.MethodDelete
	// 编号
	ID uint64 `name:"id,string" in:"path"`
}

func (req DeleteApi) Path() string {
	return "/:id"
}

func (req DeleteApi) Output(ctx context.Context) (result interface{}, err error) {
	err = modules.DeleteAPI(req.ID, storage.Database)
	return
}
