package modules

import (
	"fmt"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/models"
)

type ServerContract interface {
	GetIdentity() uint64
	SetIdentity(uint64)
	GetHost() string
	GetType() enum.ServerType
}

func NewServerContract(s models.ServerContract) ServerContract {
	if gs, ok := s.(*models.GeneralServer); ok {
		return &WebServiceServer{
			Server: Server{
				ID:         gs.ID,
				Name:       gs.Name,
				Host:       gs.Host,
				Port:       gs.Port,
				ServerType: gs.ServerType,
			},
		}
	}
	switch s.GetType() {
	case enum.SERVER_TYPE__WEB_SERVICE:
		if ws, ok := s.(*models.WebServiceServer); ok {
			return NewWebServiceServer(ws)
		}
	case enum.SERVER_TYPE__DATABASE:
		if ds, ok := s.(*models.DatabaseServer); ok {
			return NewDatabaseServer(ds)
		}
	}
	return nil
}

func NewWebServiceServer(model *models.WebServiceServer) *WebServiceServer {
	return &WebServiceServer{
		Server: Server{
			ID:         model.ID,
			Name:       model.Name,
			Host:       model.Host,
			Port:       model.Port,
			ServerType: model.ServerType,
		},
	}
}

func NewDatabaseServer(model *models.DatabaseServer) *DatabaseServer {
	return &DatabaseServer{
		Server: Server{
			ID:         model.ID,
			Name:       model.Name,
			Host:       model.Host,
			Port:       model.Port,
			ServerType: model.ServerType,
		},
		UserName: model.UserName,
		Password: model.Password,
		Extends:  model.Extends,
	}
}

type Server struct {
	// 唯一标识
	ID uint64 `json:"id,string" default:""`
	// 服务器名称
	Name string `json:"name"`
	// 地址
	Host string `json:"host"`
	// 端口
	Port uint16 `json:"port"`
	// 服务器类型
	ServerType enum.ServerType `json:"serverType"`
}

func (v *Server) SetIdentity(id uint64) {
	v.ID = id
}

func (v Server) GetIdentity() uint64 {
	return v.ID
}

func (v *Server) GetHost() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

func (v *Server) GetType() enum.ServerType {
	return v.ServerType
}

type WebServiceServer struct {
	Server
}

type DatabaseServer struct {
	Server

	// 数据库用户名
	UserName string `json:"userName"`
	// 数据库密码
	Password string `json:"password"`
	// 数据库配置扩展
	Extends map[string]string `json:"extends"`
}
