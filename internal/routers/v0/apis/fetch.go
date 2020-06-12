package apis

import (
	"context"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
)

func init() {
	Router.Register(courier.NewRouter(GetApis{}))
}

// 获取所有API
type GetApis struct {
	httpx.MethodGet
	// 分页偏移量
	Start uint64 `json:"start" name:"start" in:"query" default:"0" validate:"@uint64[0,]"`
	// 每页数据量
	Limit int64 `json:"limit" name:"limit" in:"query" default:"10" validate:"@int64[-1,100]"`
}

func (req GetApis) Path() string {
	return ""
}

type GetApisResult struct {
	NextID uint64        `json:"nextID"`
	Data   []modules.API `json:"data"`
}

func (req GetApis) Output(ctx context.Context) (result interface{}, err error) {
	resp := &GetApisResult{
		NextID: 0,
		Data:   make([]modules.API, 0),
	}

	resp.NextID, err = modules.WalkAPIs(req.Start, req.Limit, func(e storage.Element) error {
		resp.Data = append(resp.Data, *e.(*modules.API))
		return nil
	}, storage.Database)

	result = resp
	return
}
