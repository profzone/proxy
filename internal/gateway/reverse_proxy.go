package gateway

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg/pool"
	"longhorn/proxy/pkg/route"
	"time"
)

var APIServer *ReverseProxy

type ReverseProxyConf struct {
	Name            string
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}

type ReverseProxy struct {
	ReverseProxyConf
	server       *fasthttp.Server
	routes       *route.Routes
	apiContainer modules.APIContainer
}

func CreateReverseProxy(conf ReverseProxyConf) *ReverseProxy {
	return &ReverseProxy{
		ReverseProxyConf: conf,
		routes:           route.NewRoutes(),
		apiContainer:     modules.APIContainer{},
	}
}

func (s *ReverseProxy) Routes() *route.Routes {
	return s.routes
}

func (s *ReverseProxy) Start() error {
	err := s.initRoutes()
	if err != nil {
		return err
	}

	logrus.Infof("reverse proxy start listen on %s", s.ListenAddr)
	return s.startHTTP()
}

func (s *ReverseProxy) Close() error {
	return s.server.Shutdown()
}

func (s *ReverseProxy) initRoutes() error {
	_, err := modules.WalkAPIs(0, -1, func(e storage.Element) error {
		a := e.(*modules.API)
		err := s.routes.Handle(a.Method, a.URLPattern, a.ID)
		if err != nil {
			return err
		}
		s.apiContainer.Add(a)
		return nil
	}, storage.Database)
	return err
}

func (s *ReverseProxy) startHTTP() error {
	s.server = &fasthttp.Server{
		Name:            s.ReverseProxyConf.Name,
		Handler:         s.HandleHTTP,
		ErrorHandler:    s.HandleHTTPError,
		ReadTimeout:     s.ReverseProxyConf.ReadTimeout,
		WriteTimeout:    s.ReverseProxyConf.WriteTimeout,
		ReadBufferSize:  s.ReverseProxyConf.ReadBufferSize,
		WriteBufferSize: s.ReverseProxyConf.WriteBufferSize,
	}
	return s.server.ListenAndServe(s.ReverseProxyConf.ListenAddr)
}

func (s *ReverseProxy) HandleHTTP(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())
	apiID, params, _ := s.routes.Lookup(method, path)
	if apiID == 0 {
		logrus.Debugf("[%s] %s not exist", method, path)
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}
	logrus.Debugf("[%s] %s matched api: %d with params: %v", method, path, apiID, params)

	api := s.apiContainer.Get(apiID)
	if api == nil {
		ctx.Error(fmt.Sprintf("apiID %d not found", apiID), fasthttp.StatusInternalServerError)
		return
	}

	// Status check
	if api.Status == enum.API_STATUS__DOWN {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	// IP Control
	if !api.FilterIPControl(&ctx.Request) {
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
		return
	}

	// QPS
	if !api.FilterQPS() {
		ctx.Error("Too many requests", fasthttp.StatusTooManyRequests)
		return
	}

	wg := pool.WGPool.AcquireWG()
	defer pool.WGPool.ReleaseWG(wg)

	api.WalkDispatcher(func(dispatcher *modules.Dispatcher) error {
		wg.Add(1)
		go func(dispatcher *modules.Dispatcher) {
			defer wg.Done()

			resp, err := dispatcher.Dispatch(ctx, params, storage.Database)
			if err != nil {
				return
			}
			defer fasthttp.ReleaseResponse(resp)

			// TODO resp fusion
		}(dispatcher)

		return nil
	})
	wg.Wait()

	ctx.SuccessString("plain/text", "success")
}

func (s *ReverseProxy) HandleHTTPError(ctx *fasthttp.RequestCtx, err error) {

}
