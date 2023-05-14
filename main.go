package main

import (
	"AnimeLifeBackEnd/core"
	"AnimeLifeBackEnd/env"
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/middlewares"
	"AnimeLifeBackEnd/routes"
	"AnimeLifeBackEnd/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(env.Mode)
	router := gin.New()
	global.Viper = core.InitViper()
	logger := core.InitZapLogger(env.Mode)
	global.Logger = logger.Sugar()
	defer logger.Sync()
	if env.Mode == "release" {
		// register the zap logger and zap recovery which writes log to zap
		router.Use(middlewares.ZapLogger(logger))
		router.Use(middlewares.RecoveryWithZap(logger, true))
	} else {
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	}
	router.Use(middlewares.NewCors())
	// Init Mysql
	global.MysqlDB = core.InitMysqlDB()
	if global.MysqlDB != nil {
		core.RegisterTables(global.MysqlDB)
		db, _ := global.MysqlDB.DB()
		defer db.Close()
	}
	// Init Redis
	global.RedisDB = core.InitRedisDB()

	authJWT := middlewares.InitJWTAuth()
	global.WsHub = websocket.NewHub()
	go global.WsHub.Run()

	apiEndpointGroup := routes.RouterGroupApp.RouterGroup
	authEndpoints := routes.RouterGroupApp.AuthRouter
	websocketEndpoints := routes.RouterGroupApp.WebsocketRouter
	// 注册所有不需要认证的endpoint
	publicGroup := router.Group("")
	{
		authEndpoints.AddAuthRoutes(publicGroup, authJWT)
		websocketEndpoints.AddWebsocketRoutes(publicGroup)
	}
	privateGroup := router.Group("")
	privateGroup.Use(authJWT.MiddlewareFunc())
	{
		// Only use authentication middleware in release mode
		if env.Mode == "debug" {
			apiEndpointGroup.AddApiRoutes(publicGroup, publicGroup)
		} else {
			apiEndpointGroup.AddApiRoutes(privateGroup, publicGroup)
		}
	}

	router.Run(":8080")
}
