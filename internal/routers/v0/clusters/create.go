package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/http"
)

func init() {
	Router.Register(courier.NewRouter(CreateCluster{}))
}

// 创建集群
type CreateCluster struct {
	httpx.MethodPost
	Body modules.Cluster `name:"body" in:"body"`
}

func (req CreateCluster) Path() string {
	return ""
}

func (req CreateCluster) Output(ctx context.Context) (result interface{}, err error) {
	id, err := modules.CreateCluster(&req.Body, storage.Database)
	if err != nil {
		return
	}

	result = &http.IDResponse{ID: id}
	return
}
