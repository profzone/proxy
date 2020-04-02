package main

import (
	"github.com/profzone/eden-framework/pkg/application"
	"longhorn/proxy/internal/gateway"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/modules"
	"longhorn/proxy/internal/routers"
)

func main() {
	app := application.NewApplication(runner, &global.Config)
	go app.Start()
	app.WaitStop(func() error {
		err := modules.Database.Close()
		if err != nil {
			return err
		}
		return nil
	})
}

func runner(app *application.Application) error {
	// init database
	modules.Database.Init(global.Config.DBConfig, global.Config.SnowflakeConfig)

	// start administrator server
	go global.Config.GRPCServer.Serve(routers.RootRouter)
	go global.Config.HTTPServer.Serve(routers.RootRouter)

	// start gateway server
	apiServ := gateway.CreateAPIServer(gateway.APIServerConf{
		ListenAddr:      global.Config.ListenAddr,
		ReadTimeout:     global.Config.ReadTimeout,
		WriteTimeout:    global.Config.WriteTimeout,
		ReadBufferSize:  global.Config.ReadBufferSize,
		WriteBufferSize: global.Config.WriteBufferSize,
	})
	return apiServ.Start()
}
