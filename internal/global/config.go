package global

import (
	"github.com/profzone/eden-framework/pkg/courier/transport_grpc"
	"github.com/profzone/eden-framework/pkg/courier/transport_http"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/gateway"
	"time"
)

type SnowflakeConfig struct {
	Epoch      int64
	BaseNodeID int64
	NodeCount  int64
	NodeBits   uint8
	StepBits   uint8
}

type DBConfig struct {
	DBType             enum.DbType
	Endpoints          []string
	Prefix             string
	ClusterPrefix      string
	ServerPrefix       string
	BindPrefix         string
	ApiPrefix          string
	RouterPrefix       string
	OrganizationPrefix string
}

var Config = struct {
	// administrator
	GRPCServer transport_grpc.ServeGRPC
	HTTPServer transport_http.ServeHTTP

	// proxying
	APIServer       *gateway.ReverseProxy `ignored:"true"`
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int

	// db
	DBConfig

	// id generation
	SnowflakeConfig
}{
	GRPCServer: transport_grpc.ServeGRPC{
		Port: 8900,
	},
	HTTPServer: transport_http.ServeHTTP{
		Port:     8001,
		WithCORS: true,
	},

	ListenAddr:      "0.0.0.0:8000",
	ReadTimeout:     30 * time.Second,
	WriteTimeout:    60 * time.Second,
	ReadBufferSize:  0,
	WriteBufferSize: 0,

	DBConfig: DBConfig{
		DBType:             enum.DB_TYPE__ETCD,
		Endpoints:          []string{"127.0.0.1:2379"},
		Prefix:             "proxy",
		ClusterPrefix:      "proxy/clusters",
		ServerPrefix:       "proxy/servers",
		BindPrefix:         "proxy/binds",
		ApiPrefix:          "proxy/apis",
		RouterPrefix:       "proxy/routers",
		OrganizationPrefix: "proxy/organizations",
	},

	SnowflakeConfig: SnowflakeConfig{
		Epoch:      1288351723598,
		BaseNodeID: 1,
		NodeCount:  100,
		NodeBits:   10,
		StepBits:   12,
	},
}
