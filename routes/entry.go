package routes

import (
	"AnimeLifeBackEnd/routes/api"
)

type routerGroup struct {
	api.RouterGroup
	AuthRouter
	WebsocketRouter
	MetricsRouter
}

var RouterGroupApp = new(routerGroup)
