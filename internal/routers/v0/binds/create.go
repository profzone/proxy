package binds

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/http"
)

func init() {
	Router.Register(courier.NewRouter(CreateBind{}))
}

// 创建绑定
type CreateBind struct {
	httpx.MethodPost
	Body modules.Bind `name:"body" in:"body"`
}

func (req CreateBind) Path() string {
	return ""
}

func (req CreateBind) Output(ctx context.Context) (result interface{}, err error) {
	id, err := modules.CreateBind(&req.Body, storage.Database)
	if err != nil {
		return
	}

	result = &http.IDResponse{ID: id}
	return
}
