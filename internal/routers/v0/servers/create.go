package servers

import (
	"context"
	"fmt"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/http"
)

func init() {
	Router.Register(courier.NewRouter(CreateServer{}))
}

// 创建服务器
type CreateServer struct {
	httpx.MethodPost
	Body modules.ServerRequestBody `in:"body" name:"body"`
}

func (req CreateServer) Path() string {
	return ""
}

func (req CreateServer) Output(ctx context.Context) (result interface{}, err error) {
	server := req.Body.ToServer()
	if server == nil {
		return nil, fmt.Errorf("unsupport server type %v", req.Body.ServerType)
	}
	id, err := modules.CreateServer(server, storage.Database)
	if err != nil {
		return
	}

	result = &http.IDResponse{ID: id}
	return
}
