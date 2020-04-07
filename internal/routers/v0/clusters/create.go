package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
)

func init() {
	Router.Register(courier.NewRouter(CreateCluster{}))
}

// 创建集群
type CreateCluster struct {
	httpx.MethodPost
	Body modules.Cluster `name:"body" in:"body"`
}

type CreateClusterResult struct {
	// 集群ID
	ID uint64 `json:"id"`
}

func (req CreateCluster) Path() string {
	return ""
}

func (req CreateCluster) Output(ctx context.Context) (result interface{}, err error) {
	id, err := modules.CreateCluster(&req.Body, modules.Database)
	if err != nil {
		return
	}

	result = &CreateClusterResult{ID: id}
	return
}
