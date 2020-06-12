package v0

import (
	"github.com/profzone/eden-framework/pkg/courier"
	"longhorn/proxy/internal/routers/v0/apis"
	"longhorn/proxy/internal/routers/v0/binds"
	"longhorn/proxy/internal/routers/v0/clusters"
	"longhorn/proxy/internal/routers/v0/servers"
)

var Router = courier.NewRouter(V0Group{})

func init() {
	Router.Register(clusters.Router)
	Router.Register(servers.Router)
	Router.Register(binds.Router)
	Router.Register(apis.Router)
}

type V0Group struct {
	courier.EmptyOperator
}

func (V0Group) Path() string {
	return "/v0"
}
