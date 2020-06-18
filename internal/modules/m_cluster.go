package modules

import (
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/models"
)

type Cluster struct {
	// 唯一标识
	ID uint64 `json:"id,string" default:""`
	// 集群名称
	Name string `json:"name"`
	// 负载均衡类型
	LoadBalanceType enum.LoadBalanceType `json:"loadBalanceType"`
	// 负载均衡器
	loadBalancer LoadBalancer

	// 绑定的server列表
	servers []ServerContract
}

func NewCluster(model *models.Cluster) *Cluster {
	c := &Cluster{
		ID:              model.ID,
		Name:            model.Name,
		LoadBalanceType: model.LoadBalanceType,
	}

	switch model.LoadBalanceType {
	case enum.LOAD_BALANCE_TYPE__ROUND_ROBIN:
		c.loadBalancer = NewRoundRobin()
	default:
	}
	return c
}

func (v *Cluster) SetIdentity(id uint64) {
	v.ID = id
}

func (v Cluster) GetIdentity() uint64 {
	return v.ID
}

func (v *Cluster) ApplyLoadBalance(req fasthttp.Request) ServerContract {
	return v.loadBalancer.Apply(req, v.servers)
}
