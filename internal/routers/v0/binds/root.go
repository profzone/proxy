package binds

import "github.com/profzone/eden-framework/pkg/courier"

var Router = courier.NewRouter(BindGroup{})

type BindGroup struct {
	courier.EmptyOperator
}

func (BindGroup) Path() string {
	return "/binds"
}
