package main

import (
	"github.com/profzone/eden-framework/pkg/application"
	"github.com/sirupsen/logrus"
	"longhorn/proxy/internal/gateway"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/routers"
	"longhorn/proxy/internal/storage"
)

func main() {
	app := application.NewApplication(runner, &global.Config)
	go app.Start()
	app.WaitStop(func() error {
		err := storage.Database.Close()
		if err != nil {
			return err
		}
		logrus.Infof("database shutdown.")

		err = global.Config.APIServer.Close()
		if err != nil {
			return err
		}
		logrus.Infof("api server shutdown.")

		return nil
	})
}

func runner(app *application.Application) error {
	// init database
	storage.Database.Init(global.Config.DBConfig, global.Config.SnowflakeConfig)

	// start administrator server
	go global.Config.GRPCServer.Serve(routers.RootRouter)
	go global.Config.HTTPServer.Serve(routers.RootRouter)

	// start gateway server
	global.Config.APIServer = gateway.CreateAPIServer(gateway.APIServerConf{
		ListenAddr:      global.Config.ListenAddr,
		ReadTimeout:     global.Config.ReadTimeout,
		WriteTimeout:    global.Config.WriteTimeout,
		ReadBufferSize:  global.Config.ReadBufferSize,
		WriteBufferSize: global.Config.WriteBufferSize,
	})
	return global.Config.APIServer.Start()
}
