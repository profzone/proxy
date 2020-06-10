package binds

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(GetBinds{}))
}

// 分页获取绑定信息
type GetBinds struct {
	httpx.MethodGet
	// 集群ID
	ClusterID uint64 `in:"path" name:"clusterID,string"`
}

func (req GetBinds) Path() string {
	return "/:clusterID"
}

type GetBindsResult struct {
	NextID uint64         `json:"nextID"`
	Data   []modules.Bind `json:"data"`
}

func (req GetBinds) Output(ctx context.Context) (result interface{}, err error) {
	resp := &GetBindsResult{
		NextID: 0,
		Data:   make([]modules.Bind, 0),
	}

	resp.NextID, err = modules.WalkBinds(req.ClusterID, 0, -1, func(e storage.Element) error {
		resp.Data = append(resp.Data, *e.(*modules.Bind))
		return nil
	}, storage.Database)

	result = resp
	return
}
