package routes

import (
	"AnimeLifeBackEnd/websocket"
	"github.com/gin-gonic/gin"
)

type WebsocketRouter struct{}

func (r WebsocketRouter) AddWebsocketRoutes(rg *gin.RouterGroup, hub *websocket.Hub) {
	rg.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})
}
