package gateway

import (
	"github.com/valyala/fasthttp"
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
	ReverseProxyConf
}

func CreateReverseProxy(conf ReverseProxyConf) *ReverseProxy {
	return &ReverseProxy{
		ReverseProxyConf: conf,
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
	req := ctx.Request
	req.SetHost("127.0.0.1:8001")
	fasthttp.Do(&req, &ctx.Response)
}

func (s *ReverseProxy) HandleHTTPError(ctx *fasthttp.RequestCtx, err error) {

}
