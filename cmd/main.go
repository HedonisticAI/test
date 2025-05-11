package main

import (
	"test/config"
	"test/internal/app"
	"test/pkg/logger"
)

func main() {
	Logger := logger.NewLogger()
	Config := config.NewConfig()
	if Config == nil {
		Logger.Info("Bad config")
	}
	App, err := app.NewApp(*Logger, *Config)
	if err != nil {
		Logger.Error(err)
		return
	}
	App.Run()
}
