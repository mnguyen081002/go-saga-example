package main

import (
	"item-service/config"
	"item-service/internal/controller"
	"item-service/internal/http"
	"item-service/internal/lib/db"
	"item-service/internal/lib/logger"
	"item-service/internal/service"
	"item-service/utils"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func inject() fx.Option {
	return fx.Options(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		//fx.NopLogger,
		fx.Provide(
			config.NewConfig,
			utils.NewTimeoutContext,
		),
		db.GormModule,
		logger.LoggerModule,
		service.Module,
		controller.Module,
		http.Module,
	)
}

func main() {
	fx.New(inject()).Run()
}
