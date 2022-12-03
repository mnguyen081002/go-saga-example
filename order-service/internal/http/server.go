package http

import (
	"context"
	"fmt"
	"item-service/config"
	"item-service/utils/constants"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	// echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewServer)

func NewServer(lifecycle fx.Lifecycle, zap *zap.Logger, config config.Config) *gin.RouterGroup {

	switch config.Server.Env {
	case constants.Dev, constants.Local:
		gin.SetMode(gin.DebugMode)
	case constants.Prod:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	//gin.LoggerWithConfig(gin.LoggerConfig{
	//	Formatter: nil,
	//	Output:    nil,
	//	SkipPaths: nil,
	//})
	instance := gin.New()
	middleware := InitMiddleware()
	instance.Use(middleware.ErrorHandler(zap))
	instance.Use(middleware.JSONMiddleware)

	instance.Use(middleware.CORS)
	instance.Use(middleware.Logger(zap))
	//instance.Use(gozap.RecoveryWithZap(zap, true))
	instance.Use(middleware.JWT(config))

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			zap.Info("Starting HTTP server")
			go func() {
				addr := fmt.Sprint(config.Server.Host, ":", config.Server.Port)
				if err := instance.Run(addr); err != nil {
					zap.Fatal(fmt.Sprint("HTTP server failed to start %w", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			zap.Info("Stopping HTTP server")
			return nil
		},
	})

	//instance.GET("/swagger/*", echoSwagger.WrapHandler)

	return instance.Group("/api")
}
