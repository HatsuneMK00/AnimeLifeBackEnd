package main

import (
	"AnimeLifeBackEnd/core"
	"AnimeLifeBackEnd/env"
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/middlewares"
	"AnimeLifeBackEnd/routes"
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
		router.Use(middlewares.NewCors())
	}
	global.MysqlDB = core.InitMysqlDB()
	if global.MysqlDB != nil {
		core.RegisterTables(global.MysqlDB)
		db, _ := global.MysqlDB.DB()
		defer db.Close()
	}

	authJWT := middlewares.InitJWTAuth()

	apiEndpointGroup := routes.RouterGroupApp.RouterGroup
	authEndpoints := routes.RouterGroupApp.AuthRouter
	// 注册所有不需要认证的endpoint
	publicGroup := router.Group("")
	{
		authEndpoints.AddAuthRoutes(publicGroup, authJWT)
	}
	privateGroup := router.Group("")
	privateGroup.Use(authJWT.MiddlewareFunc())
	{
		apiEndpointGroup.AddApiRoutes(privateGroup, publicGroup)
	}

	router.Run(":8080")
}
