package main

import (
	"AnimeLifeBackEnd/core"
	"AnimeLifeBackEnd/env"
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/middlewares"
	"AnimeLifeBackEnd/routes"
	"AnimeLifeBackEnd/websocket"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
	"os"
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

	// Init Clerk Client
	apiKey := os.Getenv("CLERK_API_KEY")
	client, err := clerk.NewClient(apiKey)
	if err != nil {
		global.Logger.Errorf("failed to init clerk client: %v", err)
	}

	// Init JWT Auth
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
		apiEndpointGroup.AddApiRoutes(privateGroup, publicGroup)
	}
	clerkPrivateGroup := router.Group("/v2")
	clerkPublicGroup := router.Group("/v2")
	clerkPrivateGroup.Use(middlewares.NewClerkAuth(client))
	{
		apiEndpointGroup.AddApiRoutes(clerkPrivateGroup, clerkPublicGroup)
	}

	router.Run(":8080")
}
