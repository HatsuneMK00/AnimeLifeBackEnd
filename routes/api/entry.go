package api

import "github.com/gin-gonic/gin"

type RouterGroup struct {
	userRouter
	animeRecordRouter
	authCheckRouter
}

func (r RouterGroup) AddApiRoutes(rg *gin.RouterGroup, rgPublic *gin.RouterGroup) {
	apiGroup := rg.Group("/api")
	publicApiGroup := rgPublic.Group("/api")
	r.AddUserRoutes(apiGroup, publicApiGroup)
	r.AddAnimeRecordRoutes(apiGroup)
	r.AddAuthCheckRoutes(apiGroup)
}
