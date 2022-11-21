package routes

import (
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/websocket"
	"github.com/gin-gonic/gin"
)

type WebsocketRouter struct{}

func (r WebsocketRouter) AddWebsocketRoutes(rg *gin.RouterGroup) {
	rg.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(global.WsHub, c.Writer, c.Request)
	})
}
