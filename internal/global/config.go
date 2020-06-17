package global

import (
	"longhorn/proxy/internal/constants/enum"
	"time"
)

type DBConfig struct {
	DBType             enum.DbType
	Endpoints          []string
	UserName           string
	Password           string
	DatabaseName       string
	ConnectionTimeout  time.Duration
	ClusterPrefix      string
	ServerPrefix       string
	BindPrefix         string
	ApiPrefix          string
	RouterPrefix       string
	OrganizationPrefix string
}

var Config = struct {
	// proxying
	Name            string
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int

	// db
	DBConfig
}{
	Name:            "longhorn Proxy Server",
	ListenAddr:      "0.0.0.0:8000",
	ReadTimeout:     10 * time.Second,
	WriteTimeout:    30 * time.Second,
	ReadBufferSize:  0,
	WriteBufferSize: 0,

	DBConfig: DBConfig{
		DBType:             enum.DB_TYPE__MONGODB,
		Endpoints:          []string{"127.0.0.1:27017"},
		ClusterPrefix:      "clusters",
		ServerPrefix:       "servers",
		BindPrefix:         "binds",
		ApiPrefix:          "apis",
		RouterPrefix:       "routers",
		OrganizationPrefix: "organizations",
	},
}
