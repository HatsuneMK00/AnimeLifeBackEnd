package apis

import "github.com/gin-gonic/gin"

type AuthCheckApi interface {
	AuthCheck(c *gin.Context)
}

type authCheckApi struct{}

func (a authCheckApi) AuthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "authentication passed",
	})
}
