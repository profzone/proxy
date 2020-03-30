package apis

import "github.com/profzone/eden-framework/pkg/courier"

var Router = courier.NewRouter(ApiGroup{})

type ApiGroup struct {
	courier.EmptyOperator
}

func (ApiGroup) Path() string {
	return "/apis"
}
