package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/http"
)

func init() {
	Router.Register(courier.NewRouter(CreateApi{}))
}

// 创建集群
type CreateApi struct {
	httpx.MethodPost
	Body modules.API `name:"body" in:"body"`
}

func (req CreateApi) Path() string {
	return ""
}

func (req CreateApi) Output(ctx context.Context) (result interface{}, err error) {
	id, err := modules.CreateAPI(&req.Body, storage.Database)
	if err != nil {
		return
	}

	result = http.IDResponse{
		ID: id,
	}
	return
}
