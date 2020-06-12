package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(GetApi{}))
}

// 获取单个API
type GetApi struct {
	httpx.MethodGet

	ID uint64 `name:"id" in:"path"`
}

func (req GetApi) Path() string {
	return "/:id"
}

func (req GetApi) Output(ctx context.Context) (result interface{}, err error) {
	a, err := modules.GetAPI(req.ID, storage.Database)
	if err != nil {
		return
	}

	result = a
	return
}
