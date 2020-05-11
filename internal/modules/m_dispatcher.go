package modules

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/storage"
	"time"
)

type Dispatcher struct {
	Router Router `json:"router" default:""`
	// Validations
	WriteTimeout time.Duration `json:"writeTimeout" default:""`
	ReadTimeout  time.Duration `json:"readTimeout" default:""`
	ClusterID    uint64        `json:"clusterID,string"`
}

func (d *Dispatcher) Dispatch(ctx *fasthttp.RequestCtx, db storage.Storage) (*fasthttp.Response, error) {
	clusterID := d.dispatchTarget(ctx.Request)

	cluster, err := GetCluster(clusterID, db)
	if err != nil {
		return nil, err
	}

	servers := make([]*Server, 0)
	serverMap := make(map[uint64]*Server)
	_, err = WalkBinds(clusterID, 0, -1, func(e storage.Element) error {
		bind := e.(*Bind)
		server, err := GetServer(bind.ServerID, db)
		if err != nil {
			return err
		}

		servers = append(servers, server)
		serverMap[server.ID] = server
		return nil
	}, db)
	if err != nil {
		return nil, err
	}

	lb := cluster.GetLoadBalance()
	if lb == nil {
		return nil, fmt.Errorf("cluster did not set load balance type")
	}

	serverID := lb.Apply(ctx.Request, servers)
	target := serverMap[serverID]

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	err = copier.Copy(req, ctx.Request)
	if err != nil {
		return nil, err
	}
	req.SetHost(target.GetHost())

	// TODO rewrite

	err = fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *Dispatcher) dispatchTarget(originRequest fasthttp.Request) uint64 {
	return 0
}
