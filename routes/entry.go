package routes

import (
	"AnimeLifeBackEnd/routes/api"
)

type routerGroup struct {
	api.RouterGroup
	AuthRouter
}

var RouterGroupApp = new(routerGroup)
