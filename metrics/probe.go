package metrics

import "github.com/gin-gonic/gin"

func HealthProbeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
