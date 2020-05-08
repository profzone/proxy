package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(UpdateCluster{}))
}

// 更新集群
type UpdateCluster struct {
	httpx.MethodPatch
	// 编号
	ID   uint64          `name:"id,string" in:"path"`
	Body modules.Cluster `name:"body" in:"body"`
}

func (req UpdateCluster) Path() string {
	return "/:id"
}

func (req UpdateCluster) Output(ctx context.Context) (result interface{}, err error) {
	req.Body.ID = req.ID
	err = modules.UpdateCluster(&req.Body, storage.Database)
	return
}
