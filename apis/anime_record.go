package apis

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/request"
	"AnimeLifeBackEnd/global"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
)

type AnimeRecordApi interface {
	FetchAnimeRecords(c *gin.Context)
	FetchAnimeRecordsOfRating(c *gin.Context)
	FetchAnimeRecordSummary(c *gin.Context)
	AddAnimeRecord(c *gin.Context)
	UpdateAnimeRecord(c *gin.Context)
	SearchAnimeRecords(c *gin.Context)
}

type animeRecordApi struct{}

func (a animeRecordApi) FetchAnimeRecords(c *gin.Context) {
	// return 15 records from offset at a time
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
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
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
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
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
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
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
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

func (a animeRecordApi) UpdateAnimeRecord(c *gin.Context) {
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	updateRequest := request.AnimeRecordUpdateRequest{}
	err = c.ShouldBindJSON(&updateRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// update anime info if bangumiId is changed
	anime, err := animeRecordService.FetchAnimeById(updateRequest.AnimeId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to fetch anime",
		})
		return
	}
	shouldAnimeUpdate := false
	isAnimeUpdateFailed := false
	var updatedAnime entity.Anime
	if anime.BangumiId != updateRequest.BangumiId {
		shouldAnimeUpdate = true
		updatedAnime, err = animeRecordService.UpdateAnimeByBangumiId(updateRequest.BangumiId, anime)
		if err != nil {
			isAnimeUpdateFailed = true
		}
	}
	record, err := animeRecordService.UpdateAnimeRecord(updateRequest.AnimeId, userId, updateRequest.AnimeRating)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to update anime record",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": map[string]interface{}{
			"record":              record,
			"updatedAnime":        updatedAnime,
			"shouldAnimeUpdate":   shouldAnimeUpdate,
			"isAnimeUpdateFailed": isAnimeUpdateFailed,
		},
	})
}

func (a animeRecordApi) SearchAnimeRecords(c *gin.Context) {
	// return 15 records from offset at a time
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
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
	searchText, err := url.QueryUnescape(c.DefaultQuery("searchText", ""))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "searchText needs to be url encoded",
		})
		return
	}
	animeRecords, err := animeRecordService.SearchAnimeRecords(userId, offset, 15, searchText)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    animeRecords,
	})
}

func getUserIdFromJwtToken(c *gin.Context) (int, error) {
	claims := jwt.ExtractClaims(c)
	global.Logger.Debugf("claims: %v", claims)
	userId := claims[global.Config.Jwt.IdentityKey]
	global.Logger.Debugf("userId: %v", userId)
	if userId, ok := userId.(float64); ok {
		return int(userId), nil
	} else {
		err := errors.New("user id is not int")
		return 0, err
	}
}
