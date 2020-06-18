package modules

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
	"github.com/profzone/eden-framework/pkg/timelib"
	"github.com/sony/gobreaker"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/models"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/pool"
	"longhorn/proxy/pkg/route"
	"time"
)

type Dispatcher struct {
	// 路由
	Router *Router `json:"router,omitempty" default:""`
	// 熔断器
	breaker *gobreaker.CircuitBreaker

	WriteTimeout timelib.DurationString `json:"writeTimeout" default:""`
	ReadTimeout  timelib.DurationString `json:"readTimeout" default:""`
	ClusterID    uint64                 `json:"clusterID,string"`
}

func NewDispatcher(model *models.Dispatcher) *Dispatcher {
	var breaker *gobreaker.CircuitBreaker
	if model.BreakerConf != nil {
		breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:          "",
			MaxRequests:   model.BreakerConf.MaxRequests,
			Interval:      time.Duration(model.BreakerConf.Interval),
			Timeout:       time.Duration(model.BreakerConf.Timeout),
			ReadyToTrip:   BreakerStrategyTotalFailures,
			OnStateChange: breakerStateChanged,
		})
	}
	var router *Router
	if model.Router != nil {
		router = NewRouter(model.Router)
	}
	return &Dispatcher{
		Router:       router,
		breaker:      breaker,
		WriteTimeout: model.WriteTimeout,
		ReadTimeout:  model.ReadTimeout,
		ClusterID:    model.ClusterID,
	}
}

func breakerStateChanged(name string, from gobreaker.State, to gobreaker.State) {

}

func (d *Dispatcher) breakerStateChanged(name string, from gobreaker.State, to gobreaker.State) {

}

func (d *Dispatcher) Dispatch(ctx *fasthttp.RequestCtx, params route.Params, db storage.Storage) (*fasthttp.Response, error) {
	clusterID := d.dispatchTarget(&ctx.Request, params)

	var (
		cluster *Cluster
		err     error
		exist   bool
	)
	cluster, exist = ClusterContainer.GetCluster(clusterID)
	if !exist {
		model, err := models.GetCluster(clusterID, db)
		if err != nil {
			return nil, err
		}
		cluster, err = NewClusterAndServers(model)
		if err != nil {
			return nil, err
		}
		_ = ClusterContainer.AddCluster(cluster, cache.DefaultExpiration)
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

	target := cluster.ApplyLoadBalance(*req)
	if target == nil {
		return nil, fmt.Errorf("cluster did not set load balance")
	}
	req.SetHost(target.GetHost())

	cli := pool.ClientPool.AcquireClient()
	defer func() {
		pool.ClientPool.ReleaseClient(cli)
	}()
	// TODO if not set then use global config
	cli.ReadTimeout = time.Duration(d.ReadTimeout)
	cli.WriteTimeout = time.Duration(d.WriteTimeout)

	var resp *fasthttp.Response
	if d.breaker != nil {
		resp, err = d.wrapBreakerRequest(cli, req)
	} else {
		resp, err = d.forward(cli, req)
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *Dispatcher) wrapBreakerRequest(cli *fasthttp.Client, req *fasthttp.Request) (resp *fasthttp.Response, err error) {
	result, err := d.breaker.Execute(func() (resp interface{}, err error) {
		return d.forward(cli, req)
	})
	resp = result.(*fasthttp.Response)
	return
}

func (d *Dispatcher) forward(cli *fasthttp.Client, req *fasthttp.Request) (resp *fasthttp.Response, err error) {
	response := fasthttp.AcquireResponse()
	err = cli.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (d *Dispatcher) dispatchTarget(originRequest *fasthttp.Request, params route.Params) uint64 {
	if d.Router != nil && d.Router.Match(originRequest, params) && d.Router.ClusterID != 0 {
		return d.Router.ClusterID
	}
	return d.ClusterID
}
