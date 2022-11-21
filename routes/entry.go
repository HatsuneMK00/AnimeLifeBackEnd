package routes

import (
	"AnimeLifeBackEnd/routes/api"
)

type routerGroup struct {
	api.RouterGroup
	AuthRouter
	WebsocketRouter
}

var RouterGroupApp = new(routerGroup)
