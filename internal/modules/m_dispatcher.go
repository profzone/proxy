package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sony/gobreaker"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/pool"
	"longhorn/proxy/pkg/route"
	"time"
)

type Dispatcher struct {
	// 路由
	Router *Router `json:"router,omitempty" default:""`
	// 熔断器
	BreakerConf *BreakerConf `json:"breaker,omitempty" default:""`
	breaker     *gobreaker.CircuitBreaker

	WriteTimeout time.Duration `json:"writeTimeout" default:""`
	ReadTimeout  time.Duration `json:"readTimeout" default:""`
	ClusterID    uint64        `json:"clusterID,string"`
}

func (d *Dispatcher) GobDecode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(d)
	if err != nil {
		return err
	}

	if d.BreakerConf != nil {
		d.breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:          "",
			MaxRequests:   d.BreakerConf.MaxRequests,
			Interval:      d.BreakerConf.Interval,
			Timeout:       d.BreakerConf.Timeout,
			ReadyToTrip:   BreakerStrategyTotalFailures,
			OnStateChange: d.breakerStateChanged,
		})
	}

	return nil
}

func (d *Dispatcher) breakerStateChanged(name string, from gobreaker.State, to gobreaker.State) {

}

func (d *Dispatcher) Dispatch(ctx *fasthttp.RequestCtx, params route.Params, db storage.Storage) (*fasthttp.Response, error) {
	clusterID := d.dispatchTarget(&ctx.Request, params)

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

	req := fasthttp.AcquireRequest()
	defer func() {
		fasthttp.ReleaseRequest(req)
	}()

	err = copier.Copy(req, ctx.Request)
	if err != nil {
		return nil, err
	}

	if d.Router != nil && d.Router.Match(req, params) {
		err = d.Router.Rewrite(req, params)
		if err != nil {
			return nil, err
		}
	}

	lb := cluster.GetLoadBalance()
	if lb == nil {
		return nil, fmt.Errorf("cluster did not set load balance type")
	}

	serverID := lb.Apply(ctx.Request, servers)
	target := serverMap[serverID]
	req.SetHost(target.GetHost())

	cli := pool.ClientPool.AcquireClient()
	defer func() {
		pool.ClientPool.ReleaseClient(cli)
	}()
	// TODO if not set then use global config
	cli.ReadTimeout = d.ReadTimeout
	cli.WriteTimeout = d.WriteTimeout

	resp, err := d.wrapBreakerRequest(cli, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *Dispatcher) wrapBreakerRequest(cli *fasthttp.Client, req *fasthttp.Request) (resp *fasthttp.Response, err error) {
	result, err := d.breaker.Execute(func() (resp interface{}, err error) {
		response := fasthttp.AcquireResponse()
		err = cli.Do(req, response)
		if err != nil {
			return nil, err
		}
		return response, nil
	})
	resp = result.(*fasthttp.Response)
	return
}

func (d *Dispatcher) dispatchTarget(originRequest *fasthttp.Request, params route.Params) uint64 {
	if d.Router != nil && d.Router.Match(originRequest, params) && d.Router.ClusterID != 0 {
		return d.Router.ClusterID
	}
	return d.ClusterID
}
