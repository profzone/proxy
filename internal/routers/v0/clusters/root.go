package clusters

import "github.com/profzone/eden-framework/pkg/courier"

var Router = courier.NewRouter(ClusterGroup{})

type ClusterGroup struct {
	courier.EmptyOperator
}

func (ClusterGroup) Path() string {
	return "/clusters"
}
