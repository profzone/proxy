package main

import (
	"github.com/profzone/eden-framework/pkg/application"
	"github.com/sirupsen/logrus"
	"longhorn/proxy/internal/gateway"
	"longhorn/proxy/internal/global"
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

		err = gateway.APIServer.Close()
		if err != nil {
			return err
		}
		logrus.Infof("api server shutdown.")

		return nil
	})
}

func runner(app *application.Application) error {
	storage.Database.Init(global.Config.DBConfig)

	// start gateway server
	gateway.APIServer = gateway.CreateReverseProxy(gateway.ReverseProxyConf{
		Name:            global.Config.Name,
		ListenAddr:      global.Config.ListenAddr,
		ReadTimeout:     global.Config.ReadTimeout,
		WriteTimeout:    global.Config.WriteTimeout,
		ReadBufferSize:  global.Config.ReadBufferSize,
		WriteBufferSize: global.Config.WriteBufferSize,
	})
	return gateway.APIServer.Start()
}
