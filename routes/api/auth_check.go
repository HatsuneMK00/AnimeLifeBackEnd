package api

import (
	"AnimeLifeBackEnd/apis"
	"github.com/gin-gonic/gin"
)

type authCheckRouter struct{}

func (r authCheckRouter) AddAuthCheckRoutes(rg *gin.RouterGroup) {
	authCheck := rg.Group("/auth_check")
	{
		authCheck.GET("", apis.ApiGroupApp.AuthCheckApi.AuthCheck)
	}
}
