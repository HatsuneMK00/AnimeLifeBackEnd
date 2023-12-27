package middlewares

import (
	"AnimeLifeBackEnd/global"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
	"strings"
)

func NewClerkAuth(client clerk.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := c.Request.Header.Get("Authorization")
		sessionToken = strings.TrimPrefix(sessionToken, "Bearer ")

		global.Logger.Info("using clerk auth middleware")

		// verify the session
		_, err := client.VerifyToken(sessionToken)
		if err != nil {
			c.Abort()
			c.JSON(401, gin.H{
				"code":    401,
				"message": "you are not authorized",
			})
			return
		}
		tokenClaim, err := client.DecodeToken(sessionToken)
		if err != nil {
			c.Abort()
			c.JSON(401, gin.H{
				"code":    401,
				"message": "decode token error",
			})
			return
		}
		c.Set("JWT_PAYLOAD", tokenClaim.Extra)
		c.Next()
	}
}
