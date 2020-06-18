package modules

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/models"
	"longhorn/proxy/internal/storage"
)

type Cluster struct {
	// 唯一标识
	ID uint64
	// 集群名称
	Name string
	// 负载均衡类型
	LoadBalanceType enum.LoadBalanceType
	// 负载均衡器
	loadBalance LoadBalance

	// 绑定的server列表
	servers []ServerContract
}

func NewCluster(model *models.Cluster) *Cluster {
	c := &Cluster{
		ID:              model.ID,
		Name:            model.Name,
		LoadBalanceType: model.LoadBalanceType,
		servers:         make([]ServerContract, 0),
	}

	switch model.LoadBalanceType {
	case enum.LOAD_BALANCE_TYPE__ROUND_ROBIN:
		c.loadBalance = NewRoundRobin()
	default:
	}
	return c
}

func NewClusterAndServers(model *models.Cluster) (*Cluster, error) {
	cluster := NewCluster(model)
	_, err := models.WalkBinds(cluster.ID, 0, -1, func(e storage.Element) error {

		bind := e.(*models.Bind)
		model, err := models.GetServer(bind.ServerID, storage.Database)
		if err != nil {
			return err
		}
		server := NewServerContract(model)
		if server == nil {
			err := fmt.Errorf("[NewClusterAndServers] modules.NewServerContract return nil, model: %+v", model)
			logrus.Error(err)
			return err
		}
		cluster.AddServer(server)

		return nil

	}, storage.Database)

	return cluster, err
}

func (v *Cluster) SetIdentity(id uint64) {
	v.ID = id
}

func (v Cluster) GetIdentity() uint64 {
	return v.ID
}

func (v *Cluster) AddServer(s ServerContract) {
	v.servers = append(v.servers, s)
}

func (v *Cluster) ApplyLoadBalance(req fasthttp.Request) ServerContract {
	return v.loadBalance.Apply(req, v.servers)
}
