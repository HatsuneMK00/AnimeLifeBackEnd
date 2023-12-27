package api

import (
	"AnimeLifeBackEnd/apis"
	"github.com/gin-gonic/gin"
)

type userRouter struct{}

func (r userRouter) AddUserRoutes(rg *gin.RouterGroup, rgPublic *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.GET("/:id", apis.ApiGroupApp.User.FindUser)
		user.GET("/offset/:offset", apis.ApiGroupApp.User.FindUsersWithOffset)
		user.GET("/info", apis.ApiGroupApp.User.FetchUserInfo)
	}
	userPublic := rgPublic.Group("/user")
	{
		userPublic.POST("", apis.ApiGroupApp.User.AddUser)
		userPublic.POST("/webhooks/clerk_create", apis.ApiGroupApp.User.OnClerkUserCreated)
	}
}
