package binds

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(GetBind{}))
}

// 获取单个绑定信息
type GetBind struct {
	httpx.MethodGet

	ID uint64 `name:"id,string" in:"path"`
}

func (req GetBind) Path() string {
	return "/:id"
}

func (req GetBind) Output(ctx context.Context) (result interface{}, err error) {
	cluster, err := modules.GetCluster(req.ID, storage.Database)
	if err != nil {
		return
	}

	result = cluster
	return
}
