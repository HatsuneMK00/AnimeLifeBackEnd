package api

import "github.com/gin-gonic/gin"

type RouterGroup struct {
	userRouter
	animeRecordRouter
}

func (r RouterGroup) AddApiRoutes(rg *gin.RouterGroup, rgPublic *gin.RouterGroup) {
	apiGroup := rg.Group("/api")
	publicApiGroup := rgPublic.Group("/api")
	r.AddUserRoutes(apiGroup, publicApiGroup)
	r.AddAnimeRecordRoutes(apiGroup)
}
