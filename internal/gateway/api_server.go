package gateway

import (
	"github.com/valyala/fasthttp"
	"time"
)

type APIServerConf struct {
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}

type APIServer struct {
	server *fasthttp.Server
	APIServerConf
}

func CreateAPIServer(conf APIServerConf) *APIServer {
	return &APIServer{
		APIServerConf: conf,
	}
}

func (s *APIServer) Start() error {
	return s.startHTTP()
}

func (s *APIServer) Close() error {
	return s.server.Shutdown()
}

func (s *APIServer) startHTTP() error {
	s.server = &fasthttp.Server{
		Handler:         s.HandleHTTP,
		ErrorHandler:    s.HandleHTTPError,
		ReadTimeout:     s.APIServerConf.ReadTimeout,
		WriteTimeout:    s.APIServerConf.WriteTimeout,
		ReadBufferSize:  s.APIServerConf.ReadBufferSize,
		WriteBufferSize: s.APIServerConf.WriteBufferSize,
	}
	return s.server.ListenAndServe(s.APIServerConf.ListenAddr)
}

func (s *APIServer) HandleHTTP(ctx *fasthttp.RequestCtx) {
	req := ctx.Request
	req.SetHost("127.0.0.1:8001")
	fasthttp.Do(&req, &ctx.Response)
}

func (s *APIServer) HandleHTTPError(ctx *fasthttp.RequestCtx, err error) {

}
