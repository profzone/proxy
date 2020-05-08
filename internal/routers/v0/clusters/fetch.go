package clusters

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(GetClusters{}))
}

// 分页获取集群
type GetClusters struct {
	httpx.MethodGet
	// 分页偏移量
	Start uint64 `json:"start" name:"start" in:"query" default:"0" validate:"@uint64[0,]"`
	// 每页数据量
	Limit int64 `json:"limit" name:"limit" in:"query" default:"10" validate:"@int64[-1,100]"`
}

func (req GetClusters) Path() string {
	return ""
}

type GetClustersResult struct {
	NextID uint64            `json:"nextID"`
	Data   []modules.Cluster `json:"data"`
}

func (req GetClusters) Output(ctx context.Context) (result interface{}, err error) {
	resp := &GetClustersResult{
		NextID: 0,
		Data:   make([]modules.Cluster, 0),
	}

	resp.NextID, err = modules.WalkClusters(req.Start, req.Limit, func(e storage.Element) error {
		resp.Data = append(resp.Data, *e.(*modules.Cluster))
		return nil
	}, storage.Database)

	result = resp
	return
}
