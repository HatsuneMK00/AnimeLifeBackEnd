package apis

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type AnimeRecordApi interface {
	FetchAnimeRecords(c *gin.Context)
}

type animeRecordApi struct{}

func (a animeRecordApi) FetchAnimeRecords(c *gin.Context) {
	// return 15 records from offset at a time
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user id needs to be uint",
		})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "offset needs to be uint",
		})
		return
	}
	animes, err := animeRecordService.FetchAnimeRecords(uint(userId), offset, 15)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    animes,
	})
}
