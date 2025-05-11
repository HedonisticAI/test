package app

import (
	"test/config"
	create_usecase "test/internal/create/usecase"
	read_usecase "test/internal/read/usecase"
	redactor_usecase "test/internal/update_delete/usecase"
	httpserver "test/pkg/http_server"
	"test/pkg/logger"
	"test/pkg/postgres"
)

type App struct {
	Server httpserver.HttpServer
	Logger logger.Logger
}

func NewApp(Logger logger.Logger, Config config.Config) (*App, error) {
	Logger.Info("Creating new app")
	Server := httpserver.NewServer(&Config)
	DB, err := postgres.NewDB(Config)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	Logger.Debug("DB connected")
	Reader := read_usecase.NewSearchUC(DB, Logger)
	Server.MapGet("/Get", Reader.Read)
	Logger.Debug("Searching mapped")
	Creator := create_usecase.NewCreator(DB, Logger, Config.OpenApi.ApiGender, Config.OpenApi.ApiAge, Config.OpenApi.ApiNation)
	Server.MapPost("/Create", Creator.Create)
	Logger.Debug("Creating mapped")
	Redactor := redactor_usecase.NewRedactor(Logger, *DB)
	Server.MapDelete("/Delete", Redactor.Delete)
	Logger.Debug("Deleting mapped")
	Server.MapPut("/Change", Redactor.Update)
	Logger.Debug("Changing mapped")
	Logger.Info("App finished")
	return &App{Server: *Server, Logger: Logger}, nil
}

func (App *App) Run() {
	App.Logger.Info("Running App")
	App.Server.Run()
}
