package servers

import "github.com/profzone/eden-framework/pkg/courier"

var Router = courier.NewRouter(ServerGroup{})

type ServerGroup struct {
	courier.EmptyOperator
}

func (ServerGroup) Path() string {
	return "/servers"
}
