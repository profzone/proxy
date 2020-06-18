package main

import (
	"github.com/profzone/eden-framework/pkg/application"
	"github.com/profzone/eden-framework/pkg/context"
	"longhorn/proxy/internal/gateway"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

func main() {
	app := application.NewApplication(runner, &global.Config)
	go app.Start()
	app.WaitStop(func(ctx *context.WaitStopContext) error {
		ctx.Cancel()
		return nil
	})
}

func runner(app *application.Application) error {
	storage.Database.Init(global.Config.DBConfig, app.Context())

	// start gateway server
	gateway.APIServer = gateway.CreateReverseProxy(gateway.ReverseProxyConf{
		Name:            global.Config.Name,
		ListenAddr:      global.Config.ListenAddr,
		ReadTimeout:     global.Config.ReadTimeout,
		WriteTimeout:    global.Config.WriteTimeout,
		ReadBufferSize:  global.Config.ReadBufferSize,
		WriteBufferSize: global.Config.WriteBufferSize,
	}, app.Context())
	return gateway.APIServer.Start()
}
