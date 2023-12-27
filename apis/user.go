package apis

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/response"
	"AnimeLifeBackEnd/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserApi interface {
	FindUser(c *gin.Context)
	FindUsersWithOffset(c *gin.Context)
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	FetchUserInfo(c *gin.Context)
	OnClerkUserCreated(c *gin.Context)
}

type userApi struct{}

func (api userApi) FindUser(c *gin.Context) {
	userId := c.Param("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "user id needs to be uint",
		})
		return
	}
	if user, ok := userService.FindUser(uint(id)); ok {
		c.JSON(http.StatusOK, response.Response{
			Code: http.StatusOK,
			Data: user,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code:    404,
			Message: "user not found",
		})
	}
}

func (api userApi) FindUsersWithOffset(c *gin.Context) {
	param := c.Param("offset")
	offset, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "user id needs to be uint",
		})
		return
	}
	if users, ok := userService.FindUsersWithOffset(offset); ok {
		c.JSON(http.StatusOK, response.Response{
			Code: http.StatusOK,
			Data: users,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code:    404,
			Message: "user not found",
		})
	}
}

func (api userApi) AddUser(c *gin.Context) {
	var user entity.User
	var result *entity.User
	var rowAffected int64
	if err := c.ShouldBindJSON(&user); err == nil {
		result, rowAffected = userService.AddUser(&user)
	}
	if rowAffected > 0 {
		c.JSON(http.StatusOK, response.Response{
			Code:    http.StatusOK,
			Message: "Add user success",
		})
		global.Logger.Infof("add user: %v", result.Username)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "unable to add user"})
	}
}

func (api userApi) OnClerkUserCreated(c *gin.Context) {
	var clerkResponse map[string]interface{}
	/*
			Clerk Data has the following format
			{
				"data": {
					"object": "event",
		  			"type": "user.created"
					"id": "user_apsiodfjaposdifjaopij",
					"username": null,
					"email_addresses": [
					  {
						"email_address": "example@example.org",
						"id": "idn_29w83yL7CwVlJXylYLxcslromF1",
						"linked_to": [],
						"object": "email_address",
						"verification": {
						  "status": "verified",
						  "strategy": "ticket"
						}
					  }
					],
				}
			}
	*/
	if err := c.ShouldBindJSON(&clerkResponse); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to bind json",
		})
		return
	}
	// Check request body
	object := clerkResponse["object"].(string)
	eventType := clerkResponse["type"].(string)
	if object != "event" || eventType != "user.created" {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "Fail to add user",
		})
		global.Logger.Errorf("Fail to add user, clerk data is not user.created event")
	}

	clerkData := clerkResponse["data"].(map[string]interface{})
	clerkId := clerkData["id"].(string)
	username := clerkData["username"].(string)
	// Leave email address unimplemented first

	user := entity.User{
		ClerkId:  clerkId,
		Username: username,
	}
	result, rowAffected := userService.AddUser(&user)
	if rowAffected > 0 {
		c.JSON(http.StatusOK, response.Response{
			Code:    http.StatusOK,
			Message: "Add user success",
		})
		global.Logger.Infof("add user: %v", result.Username)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "unable to add user"})
	}
}

func (api userApi) FetchUserInfo(c *gin.Context) {
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	if user, ok := userService.FindUser(uint(userId)); ok {
		c.JSON(http.StatusOK, response.Response{
			Code: http.StatusOK,
			Data: user,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code:    404,
			Message: "user not found",
		})
	}
}

func (api userApi) UpdateUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (api userApi) DeleteUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
