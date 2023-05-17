package routes

import (
	"AnimeLifeBackEnd/apis"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (r AuthRouter) AddAuthRoutes(rg *gin.RouterGroup, authJWT *jwt.GinJWTMiddleware) {
	rg.POST("/login", authJWT.LoginHandler)
	rg.GET("/refresh_token", authJWT.RefreshHandler)
	rg.GET("/login_via_email", apis.ApiGroupApp.AuthApi.LoginViaEmail)
	//rg.POST("/register", apis.ApiGroupApp.AuthApi.RegisterViaEmail)
}
