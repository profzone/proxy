package global

import (
	"github.com/profzone/eden-framework/pkg/courier/transport_grpc"
	"github.com/profzone/eden-framework/pkg/courier/transport_http"
	"time"
)

var Config = struct {
	// administrator
	GRPCServer transport_grpc.ServeGRPC
	HTTPServer transport_http.ServeHTTP

	// proxying
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}{
	GRPCServer: transport_grpc.ServeGRPC{
		Port: 8900,
	},
	HTTPServer: transport_http.ServeHTTP{
		Port:     8000,
		WithCORS: true,
	},

	ListenAddr:      "0.0.0.0:8000",
	ReadTimeout:     30 * time.Second,
	WriteTimeout:    60 * time.Second,
	ReadBufferSize:  0,
	WriteBufferSize: 0,
}
