package api

import "github.com/gin-gonic/gin"

type RouterGroup struct {
	userRouter
	animeRecordRouter
	authCheckRouter
}

func (r RouterGroup) AddApiRoutes(rg *gin.RouterGroup) {
	apiGroup := rg.Group("/api")
	r.AddUserRoutes(apiGroup)
	r.AddAnimeRecordRoutes(apiGroup)
	r.AddAuthCheckRoutes(apiGroup)
}
