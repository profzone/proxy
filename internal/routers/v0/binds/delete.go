package binds

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(DeleteBind{}))
}

// 解除绑定
type DeleteBind struct {
	httpx.MethodDelete
	// 集群ID
	ClusterID uint64 `in:"path" name:"clusterID,string"`
	// 服务器ID
	ServerID uint64 `in:"path" name:"serverID,string,omitempty"`
}

func (req DeleteBind) Path() string {
	return "/:clusterID/servers/:serverID"
}

func (req DeleteBind) Output(ctx context.Context) (result interface{}, err error) {
	err = modules.DeleteBind(req.ClusterID, req.ServerID, storage.Database)
	return
}
