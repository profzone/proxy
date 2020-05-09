package gateway

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/pkg/route"
	"time"
)

type ReverseProxyConf struct {
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}

type ReverseProxy struct {
	server *fasthttp.Server
	routes *route.Routes
	ReverseProxyConf
}

func CreateReverseProxy(conf ReverseProxyConf) *ReverseProxy {
	return &ReverseProxy{
		ReverseProxyConf: conf,
		routes:           route.NewRoutes(),
	}
}

func (s *ReverseProxy) Start() error {
	return s.startHTTP()
}

func (s *ReverseProxy) Close() error {
	return s.server.Shutdown()
}

func (s *ReverseProxy) startHTTP() error {
	s.server = &fasthttp.Server{
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
	logrus.Debug(method, path)
	apiID, params, exist := s.routes.Lookup(method, path)
	if !exist {
		logrus.Errorf("[%s] %s not exist", method, path)
	}
	logrus.Infof("[%s] %s matched api: %d with params: %v", method, path, apiID, params)
}

func (s *ReverseProxy) HandleHTTPError(ctx *fasthttp.RequestCtx, err error) {

}
