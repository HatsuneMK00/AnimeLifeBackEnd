package apis

import (
	"AnimeLifeBackEnd/entity/request"
	"AnimeLifeBackEnd/global"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AnimeRecordApi interface {
	FetchAnimeRecords(c *gin.Context)
	FetchAnimeRecordsOfRating(c *gin.Context)
	FetchAnimeRecordSummary(c *gin.Context)
	AddAnimeRecord(c *gin.Context)
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

func (a animeRecordApi) FetchAnimeRecordsOfRating(c *gin.Context) {
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
	rating, err := strconv.Atoi(c.Param("rating"))
	if err != nil || rating < -1 || rating > 4 {
		c.JSON(400, gin.H{
			"message": "rating needs to be uint and between [-1, 4]",
		})
		return
	}
	animes, err := animeRecordService.FetchAnimeRecordsOfRating(uint(userId), offset, 15, rating)
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

func (a animeRecordApi) FetchAnimeRecordSummary(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user id needs to be uint",
		})
		return
	}
	summary, err := animeRecordService.FetchAnimeRecordSummary(uint(userId))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    summary,
	})
}

func (a animeRecordApi) AddAnimeRecord(c *gin.Context) {
	//	read request body from json
	var record request.AnimeRecordRequest
	err := c.ShouldBindJSON(&record)
	if err != nil {
		global.Logger.Error(err)
		c.JSON(400, gin.H{
			"message": "invalid request body",
		})
		return
	}
	//	validate request body
	if !(record.AnimeRating == -1 || // 想看
		record.AnimeRating == 1 || // 一般
		record.AnimeRating == 2 || // 不错
		record.AnimeRating == 3 || // 好看
		record.AnimeRating == 4) { // 神作
		c.JSON(400, gin.H{
			"message": "invalid anime rating",
		})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user id needs to be uint",
		})
		return
	}

	anime, err := animeRecordService.AddNewAnime(record.AnimeName)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to add anime",
		})
		return
	}
	animeRecord, err := animeRecordService.AddNewAnimeRecord(int(anime.ID), userId, record.AnimeRating)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to add anime record, maybe the record already exists",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": map[string]interface{}{
			"anime":  anime,
			"record": animeRecord,
		},
	})
}
